package configurator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"
	"time"

	autoscale "github.com/ferretcode/switchyard/autoscale/pkg/types"
	"github.com/ferretcode/switchyard/configurator/internal/railway"
	"github.com/ferretcode/switchyard/configurator/internal/railway/gql"
	"github.com/ferretcode/switchyard/configurator/internal/types"
	featureflags "github.com/ferretcode/switchyard/feature-flags/pkg/types"
	incident "github.com/ferretcode/switchyard/incident/pkg/types"
	scheduler "github.com/ferretcode/switchyard/scheduler/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var ConfigRegistry = map[string]any{
	"scheduler":     scheduler.Config{},
	"autoscale":     autoscale.Config{},
	"feature-flags": featureflags.Config{},
	"incident":      incident.Config{},
	"locomotive":    LocomotiveConfig{},
}

type ConfiguratorService struct {
	Logger    *slog.Logger
	Config    *types.Config
	GqlClient *railway.GraphQLClient
	Context   context.Context
}

func NewConfiguratorService(logger *slog.Logger, config *types.Config, gqlClient *railway.GraphQLClient, context context.Context) ConfiguratorService {
	return ConfiguratorService{
		Logger:    logger,
		Config:    config,
		GqlClient: gqlClient,
		Context:   context,
	}
}

func (c *ConfiguratorService) UpdateConfig(w http.ResponseWriter, r *http.Request) error {
	serviceName := chi.URLParam(r, "service")
	schema, ok := ConfigRegistry[serviceName]
	if !ok {
		http.Error(w, "unknown service", http.StatusNotFound)
		return nil
	}

	targetServiceId, err := c.findServiceIdFromServiceName(serviceName)
	if err != nil {
		return fmt.Errorf("error finding service from service name: %w", err)
	}

	cfgPtr := reflect.New(reflect.TypeOf(schema)).Interface()

	if err := json.NewDecoder(r.Body).Decode(cfgPtr); err != nil {
		return fmt.Errorf("error decoding config schema: %w", err)
	}

	if err := validator.New().Struct(cfgPtr); err != nil {
		return fmt.Errorf("error validating config schema: %w", err)
	}

	envMap, err := structToEnvMap(cfgPtr)
	if err != nil {
		return fmt.Errorf("error converting config schema to environment variables: %w", err)
	}

	err = c.variableCollectionUpsert(targetServiceId, envMap)
	if err != nil {
		return fmt.Errorf("error upserting service environment variables: %w", err)
	}

	w.WriteHeader(200)
	return nil
}

func (c *ConfiguratorService) findServiceIdFromServiceName(serviceName string) (string, error) {
	responseBytes, err := c.GqlClient.Client.ExecRaw(c.Context, gql.ProjectQuery, map[string]any{
		"id": c.Config.RailwayProjectId,
	})
	if err != nil {
		return "", err
	}

	projectData := gql.ProjectData{}

	if err := json.Unmarshal(responseBytes, &projectData); err != nil {
		return "", err
	}

	for _, service := range projectData.Project.Services.Edges {
		if service.Node.Name == serviceName {
			return service.Node.Id, nil
		}
	}

	return "", errors.New("that service does not exist")
}

func (c *ConfiguratorService) variableCollectionUpsert(serviceId string, vars map[string]string) error {
	variables := map[string]interface{}{
		"environmentId": c.Config.RailwayEnvironmentId,
		"projectId":     c.Config.RailwayProjectId,
		"serviceId":     serviceId,
		"variables":     vars,
	}

	_, err := c.GqlClient.Client.ExecRaw(c.Context, gql.VariableCollectionUpsertQuery, variables)
	return err
}

func structToEnvMap(cfg interface{}) (map[string]string, error) {
	result := make(map[string]string)

	val := reflect.ValueOf(cfg)
	typ := reflect.TypeOf(cfg)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, errors.New("structToEnvMap: input must be a struct or pointer to struct")
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		if field.PkgPath != "" {
			continue
		}

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		if isZeroValue(value) {
			continue
		}

		var strValue string

		switch value.Kind() {
		case reflect.String:
			strValue = value.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value.Type() == reflect.TypeOf(time.Duration(0)) {
				strValue = value.Interface().(time.Duration).String()
			} else {
				strValue = strconv.FormatInt(value.Int(), 10)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			strValue = strconv.FormatUint(value.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			strValue = strconv.FormatFloat(value.Float(), 'f', -1, 64)
		case reflect.Bool:
			strValue = strconv.FormatBool(value.Bool())
		default:
			strValue = fmt.Sprintf("%v", value.Interface())
		}

		result[envTag] = strValue
	}

	return result, nil
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Array:
		return v.Len() == 0
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Duration(0)) {
			return v.Interface().(time.Duration) == 0
		}
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}
