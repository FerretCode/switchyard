package flags

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ferretcode/switchyard/feature-flags/internal/repositories"
	"github.com/ferretcode/switchyard/feature-flags/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
)

type FlagsService struct {
	Logger  *slog.Logger
	Config  *types.Config
	Queries *repositories.Queries
	Context context.Context
}

func NewFlagsService(logger *slog.Logger, config *types.Config, queries *repositories.Queries, context context.Context) FlagsService {
	return FlagsService{
		Logger:  logger,
		Config:  config,
		Context: context,
		Queries: queries,
	}
}

func (f *FlagsService) Create(w http.ResponseWriter, r *http.Request) error {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	requestFlag := types.Flag{}

	if err := json.Unmarshal(bytes, &requestFlag); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}

	if requestFlag.Name == "" {
		http.Error(w, "feature flag name must not be empty", http.StatusBadRequest)
		return nil
	}

	featureFlag, err := f.Queries.CreateFeatureFlag(f.Context, repositories.CreateFeatureFlagParams{
		Name:    requestFlag.Name,
		Enabled: requestFlag.Enabled,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == PostgresDuplicateKeyErrorCode {
				http.Error(w, "a feature flag with that name already exists", http.StatusConflict)

				return nil
			}
		} else {
			return fmt.Errorf("error creating feature flag: %w", err)
		}
	}

	fields, operators, values := transformRules(requestFlag.Rules)

	rules, err := f.Queries.BulkCreateRulesForFeatureFlag(f.Context, repositories.BulkCreateRulesForFeatureFlagParams{
		Column1: featureFlag.ID,
		Column2: fields,
		Column3: operators,
		Column4: values,
	})
	if err != nil {
		return fmt.Errorf("error creating feature flag rules: %w", err)
	}

	ruleIds := getRuleIdsFromRules(rules)

	if err := f.Queries.BulkAssociateFeatureFlagWithRules(f.Context, repositories.BulkAssociateFeatureFlagWithRulesParams{
		FeatureFlagID: featureFlag.ID,
		Column2:       ruleIds,
	}); err != nil {
		return fmt.Errorf("error creating feature flag to rules mappings: %w", err)
	}

	w.WriteHeader(200)
	return nil
}

func (f *FlagsService) Get(w http.ResponseWriter, r *http.Request) error {
	flagName := chi.URLParam(r, "name")

	rows, err := f.Queries.GetFeatureFlagByNameWithRules(f.Context, flagName)
	if err != nil {
		return fmt.Errorf("failed to fetch feature flag: %w", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("feature flag with name %s not found", flagName)
	}

	flag := rows[0]

	getFlagResponse := GetFlagResponse{
		Name:    flag.FeatureFlagName,
		Enabled: flag.FeatureFlagEnabled,
		Rules:   []types.Rule{},
	}

	for _, row := range rows {
		if !row.RuleID.Valid {
			continue
		}

		rule := types.Rule{
			Field:    row.RuleField.String,
			Operator: row.RuleOperator.String,
			Value:    row.RuleValue.String,
		}
		getFlagResponse.Rules = append(getFlagResponse.Rules, rule)
	}

	responseBytes, err := json.Marshal(getFlagResponse)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(responseBytes)

	return nil
}

func (f *FlagsService) Update(w http.ResponseWriter, r *http.Request) error {
	flagName := chi.URLParam(r, "name")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	updatedFlag := types.Flag{}

	if err := json.Unmarshal(bytes, &updatedFlag); err != nil {
		return fmt.Errorf("error parsing updates: %w", err)
	}

	fields, operators, values := transformRules(updatedFlag.Rules)

	_, err = f.Queries.UpsertFeatureFlagByNameWithRules(f.Context, repositories.UpsertFeatureFlagByNameWithRulesParams{
		Name:    flagName,
		Enabled: updatedFlag.Enabled,
		Column3: fields,
		Column4: operators,
		Column5: values,
	})
	if err != nil {
		return fmt.Errorf("error updating feature flag: %w", err)
	}

	w.WriteHeader(200)

	return nil
}

func (f *FlagsService) Delete(w http.ResponseWriter, r *http.Request) error {
	flagName := chi.URLParam(r, "name")

	flag, err := f.Queries.GetFeatureFlagByName(f.Context, flagName)
	if err != nil {
		return fmt.Errorf("error fetching feature flag: %w", err)
	}

	err = f.Queries.DeleteRulesByFeatureFlag(f.Context, flag.ID)
	if err != nil {
		return fmt.Errorf("error deleting feature flag rules: %w", err)
	}

	err = f.Queries.DeleteFeatureFlag(f.Context, flag.ID)
	if err != nil {
		return fmt.Errorf("error deleting feature flag: %w", err)
	}

	w.WriteHeader(200)

	return nil
}

func (f *FlagsService) Evaluate(w http.ResponseWriter, r *http.Request) error {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	evaluateUserContext := EvaluateUserContext{}

	if err := json.Unmarshal(bytes, &evaluateUserContext); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}

	flagName := chi.URLParam(r, "name")

	rows, err := f.Queries.GetFeatureFlagByNameWithRules(f.Context, flagName)
	if err != nil {
		return fmt.Errorf("failed to fetch feature flag: %w", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("feature flag with name %s not found", flagName)
	}

	flag := rows[0]

	rules := []types.Rule{}

	for _, row := range rows {
		if !row.RuleID.Valid {
			continue
		}

		rule := types.Rule{
			Field:    row.RuleField.String,
			Operator: row.RuleOperator.String,
			Value:    row.RuleValue.String,
		}
		rules = append(rules, rule)
	}

	f.Logger.Info("rules", "rules", rules)

	evaluateFlagResponse := EvaluateFlagResponse{}

	evaluateFlagResponse.EnabledForUser = evaluateFlag(types.Flag{
		Enabled: flag.FeatureFlagEnabled,
		Rules:   rules,
	}, evaluateUserContext.UserContext)

	responseBytes, err := json.Marshal(evaluateFlagResponse)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(responseBytes)

	return nil

}

func transformRules(rules []types.Rule) ([]string, []string, []string) {
	fields := make([]string, len(rules))
	operators := make([]string, len(rules))
	values := make([]string, len(rules))

	for i, rule := range rules {
		fields[i] = rule.Field
		operators[i] = rule.Operator
		values[i] = rule.Value
	}

	return fields, operators, values
}

func getRuleIdsFromRules(rules []repositories.Rule) []int32 {
	ids := make([]int32, len(rules))

	for i, rule := range rules {
		ids[i] = rule.ID
	}

	return ids
}
