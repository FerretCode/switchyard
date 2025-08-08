package railway

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ferretcode/switchyard/autoscale/internal/railway/gql"
	"github.com/ferretcode/switchyard/autoscale/pkg/types"
)

type QueryService struct {
	gqlClient *GraphQLClient
	ctx       context.Context
	config    types.Config
	logger    *slog.Logger
}

func NewQueryService(gqlClient *GraphQLClient, ctx context.Context, config types.Config, logger *slog.Logger) QueryService {
	return QueryService{
		gqlClient: gqlClient,
		ctx:       ctx,
		config:    config,
		logger:    logger,
	}
}

func (q *QueryService) QueryProjectData() (*gql.ProjectData, error) {
	project := gql.ProjectData{}

	response, err := q.gqlClient.Client.ExecRaw(q.ctx, gql.ProjectQuery, map[string]any{
		"id": q.config.RailwayProjectId,
	})
	if err != nil {
		q.logger.Error("error executing graphql request", "err", err)
		return nil, err
	}

	if err := json.Unmarshal(response, &project); err != nil {
		q.logger.Error("error unmarshalling response bytes into project struct", "err", err)
		return nil, err
	}

	return &project, nil
}

func (q *QueryService) QueryServiceMetrics(serviceId string, startDate string) (*gql.MetricsData, error) {
	response, err := q.gqlClient.Client.ExecRaw(q.ctx, gql.MetricsQuery, map[string]any{
		"serviceId":    serviceId,
		"measurements": []string{"CPU_USAGE", "MEMORY_USAGE_GB", "CPU_LIMIT", "MEMORY_LIMIT_GB"},
		"startDate":    startDate,
	})
	if err != nil {
		q.logger.Error("error executing graphql request", "err", err)
		return nil, err
	}

	metrics := gql.MetricsData{}

	if err := json.Unmarshal(response, &metrics); err != nil {
		q.logger.Error("error unmarshalling repsonse bytes into metrics struct", "err", err)
		return nil, err
	}

	return &metrics, nil
}

func (q *QueryService) MutationUpdateReplicas(environmentId string, serviceId string, regionName string, newReplicas int) error {
	vars := map[string]interface{}{
		"environmentId": environmentId,
		"serviceId":     serviceId,
		"multiRegionConfig": map[string]interface{}{
			regionName: map[string]interface{}{
				"numReplicas": newReplicas,
			},
		},
	}

	q.logger.Info("making autoscale mutation request", "variables", vars)

	_, err := q.gqlClient.Client.ExecRaw(q.ctx, gql.UpdateRegionsQuery, vars)
	return err
}

func (q *QueryService) MutationServiceInstanceRedeploy(environmentId string, serviceId string) error {
	_, err := q.gqlClient.Client.ExecRaw(q.ctx, gql.ServiceInstanceDeployQuery, map[string]interface{}{
		"environmentId": environmentId,
		"serviceId":     serviceId,
	})
	return err
}
