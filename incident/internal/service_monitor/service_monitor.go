package servicemonitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/ferretcode/switchyard/incident/internal/railway"
	"github.com/ferretcode/switchyard/incident/internal/railway/gql"
	"github.com/ferretcode/switchyard/incident/internal/types"
	"github.com/ferretcode/switchyard/incident/internal/webhook"
)

type ServiceMonitorService struct {
	Logger             *slog.Logger
	Config             *types.Config
	IncidentStats      *types.IncidentStats
	PrometheusCounters *types.PrometheusCounters
	GraphQLClient      *railway.GraphQLClient
	DeploymentCache    *types.DeploymentCache
	WebhookService     *webhook.WebhookService
}

func NewServiceMonitorService(logger *slog.Logger, incidentStats *types.IncidentStats, config *types.Config, prometheusCounters *types.PrometheusCounters, graphQLClient *railway.GraphQLClient, deploymentCache *types.DeploymentCache, webhookService *webhook.WebhookService) ServiceMonitorService {
	return ServiceMonitorService{
		Logger:             logger,
		IncidentStats:      incidentStats,
		Config:             config,
		PrometheusCounters: prometheusCounters,
		GraphQLClient:      graphQLClient,
		DeploymentCache:    deploymentCache,
		WebhookService:     webhookService,
	}
}

func (s *ServiceMonitorService) MonitorServices(done chan bool) {
	ticker := time.NewTicker(time.Duration(s.Config.ServiceMonitorPollingRate) * time.Second)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			for _, serviceId := range s.Config.RailwayServiceIds {
				err := s.pollServiceId(serviceId)
				if err != nil {
					s.Logger.Error("error polling service", "service-id", serviceId, "err", err)
				}
			}
		}
	}
}

func (s *ServiceMonitorService) StartCacheCleaningJob(done chan bool) {
	ticker := time.NewTicker(time.Minute * 10)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			s.DeploymentCache.Mutex.Lock()

			for _, serviceId := range s.Config.RailwayServiceIds {
				serviceData, err := s.getServiceDataForServiceId(serviceId)
				if err != nil {
					s.Logger.Error("error fetching service data while cleaning deployment cache", "service-id", serviceId, "err", err)
				}

				for deploymentId := range s.DeploymentCache.Deployments {
					deploymentExists := slices.Contains(serviceData.Service.Deployments.Edges, gql.Edge{Node: gql.Node{Id: deploymentId}})

					if !deploymentExists {
						delete(s.DeploymentCache.Deployments, deploymentId)
					}
				}
			}

			s.DeploymentCache.Mutex.Unlock()
		}
	}
}

func (s *ServiceMonitorService) pollServiceId(serviceId string) error {
	serviceData, err := s.getServiceDataForServiceId(serviceId)
	if err != nil {
		return err
	}

	s.DeploymentCache.Mutex.Lock()
	defer s.DeploymentCache.Mutex.Unlock()

	for _, deployment := range serviceData.Service.Deployments.Edges {
		lastStatus, ok := s.DeploymentCache.Deployments[deployment.Node.Id]
		if !ok {
			s.DeploymentCache.Deployments[deployment.Node.Id] = deployment.Node.Status
			lastStatus = deployment.Node.Status
		}

		// a non-present key indicates that the deployment is new.
		// we are still interested in sending status updates if
		// the consuming service is interested in them, regardless
		// of the deployment's age
		if lastStatus != deployment.Node.Status || !ok {
			interestingStatus := slices.Contains(s.Config.ServiceMonitorInterestedStatusChanges, deployment.Node.Status)
			if !interestingStatus {
				continue
			}

			s.Logger.Info("deployment status change", "last-status", lastStatus, "new-status", deployment.Node.Status)

			err := s.WebhookService.SendDeploymentIncidentReport(
				fmt.Sprintf("deployment status changed: %s -> %s", lastStatus, deployment.Node.Status),
				serviceData.Service.Id,
				deployment.Node.Id,
			)
			if err != nil {
				return err
			}
		}

		s.DeploymentCache.Deployments[deployment.Node.Id] = deployment.Node.Status
	}

	return nil
}

func (s *ServiceMonitorService) getServiceDataForServiceId(serviceId string) (gql.ServiceData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.Config.ServiceMonitorPollingTimeout)*time.Second)
	defer cancel()

	res, err := s.GraphQLClient.Client.ExecRaw(ctx, gql.ServiceQuery, map[string]any{
		"serviceId": serviceId,
	})
	if err != nil {
		return gql.ServiceData{}, err
	}

	serviceData := gql.ServiceData{}

	if err := json.Unmarshal(res, &serviceData); err != nil {
		return gql.ServiceData{}, err
	}

	return serviceData, nil
}
