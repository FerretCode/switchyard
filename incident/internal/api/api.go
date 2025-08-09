package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ferretcode/switchyard/incident/internal/repositories"
	"github.com/ferretcode/switchyard/incident/pkg/types"
)

type ApiService struct {
	Logger  *slog.Logger
	Config  *types.Config
	Queries *repositories.Queries
	Context context.Context
}

func NewApiService(logger *slog.Logger, config *types.Config, queries *repositories.Queries, context context.Context) ApiService {
	return ApiService{
		Logger:  logger,
		Config:  config,
		Queries: queries,
		Context: context,
	}
}

func (a *ApiService) ListIncidentReports(w http.ResponseWriter, r *http.Request) error {
	deploymentRows, err := a.Queries.ListIncidentReportsWithServiceID(a.Context, 50)
	if err != nil {
		return fmt.Errorf("error listing deployment reports: %w", err)
	}

	generalRows, err := a.Queries.ListIncidentReportsWithoutServiceID(a.Context, 50)
	if err != nil {
		return fmt.Errorf("error listing general incident reports: %w", err)
	}

	var deploymentReports []types.IncidentReport
	var generalReports []types.IncidentReport

	for _, row := range deploymentRows {
		incidentReport := types.IncidentReport{
			ServiceId:     row.ServiceID,
			DeploymentId:  row.DeploymentID,
			EnvironmentId: row.EnvironmentID,
			Message:       row.Message,
			Timestamp:     row.Timestamp,
		}

		deploymentReports = append(deploymentReports, incidentReport)
	}

	for _, row := range generalRows {
		incidentReport := types.IncidentReport{
			Message:   row.Message,
			Timestamp: row.Timestamp,
		}

		generalReports = append(generalReports, incidentReport)
	}

	response := ListIncidentReportsResponse{
		DeploymentReports: deploymentReports,
		GeneralReports:    generalReports,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshalling response: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(responseBytes)

	return nil
}
