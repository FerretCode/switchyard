package autoscale

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ferretcode/switchyard/autoscale/internal/repositories"
	"github.com/go-chi/chi/v5"
)

func (a *AutoscaleService) RegisterService(w http.ResponseWriter, r *http.Request) error {
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading reques body: %w", err)
	}

	registerServiceRequest := RegisterServiceRequest{}

	if err := json.Unmarshal(requestBytes, &registerServiceRequest); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}

	if registerServiceRequest.ServiceId == "" {
		http.Error(w, "service id must not be empty", http.StatusBadRequest)
		return nil
	}

	_, err = a.Queries.CreateService(a.Context, repositories.CreateServiceParams{
		ServiceID: registerServiceRequest.ServiceId,
		JobName: sql.NullString{
			String: registerServiceRequest.JobName,
			Valid:  true,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating service: %w", err)
	}

	w.WriteHeader(200)
	return nil
}

func (a *AutoscaleService) UnregisterService(w http.ResponseWriter, r *http.Request) error {
	serviceId := chi.URLParam(r, "id")

	err := a.Queries.DeleteService(a.Context, serviceId)
	if err != nil {
		return fmt.Errorf("error deleting service: %w", err)
	}

	w.WriteHeader(200)
	return nil
}

func (a *AutoscaleService) ListServices(w http.ResponseWriter, r *http.Request) error {
	project, err := a.GqlQueries.QueryProjectData()
	if err != nil {
		return fmt.Errorf("error querying project data: %w", err)
	}

	currentEnvironmentName := ""
	for _, environment := range project.Project.Environments.Edges {
		if environment.Node.Id == a.Config.RailwayEnvironmentId {
			currentEnvironmentName = environment.Node.Name
			break
		}
	}

	var serviceContexts []ServiceContext

	for _, service := range project.Project.Services.Edges {
		serviceContext := ServiceContext{}

		currentReplicas := a.getCurrentReplicas(project, service.Node.Id)

		serviceContext.ServiceId = service.Node.Id
		serviceContext.ProjectId = project.Project.Id
		serviceContext.EnvironmentId = a.Config.RailwayEnvironmentId
		serviceContext.EnvironmentName = currentEnvironmentName
		serviceContext.Replicas = currentReplicas
		serviceContext.MinReplicas = a.Config.MinReplicaCount
		serviceContext.MaxReplicas = a.Config.MaxReplicaCount
		serviceContext.CpuUpscaleThreshold = a.Config.RailwayCpuUpscaleThreshold * 100
		serviceContext.MemoryUpscaleThreshold = a.Config.RailwayMemoryUpscaleThreshold * 100
		serviceContext.UpscaleCooldown = a.Config.UpscaleCooldown.String()
		serviceContext.DownscaleCooldown = a.Config.DownscaleCooldown.String()
		serviceContext.LastScaledAt = service.Node.ServiceInstances.Edges[0].Node.LatestDeployment.CreatedAt

		dbService, err := a.Queries.GetService(a.Context, service.Node.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				serviceContext.JobName = service.Node.Name
				serviceContext.Enabled = false
			} else {
				return fmt.Errorf("error fetching service from database: %w", err)
			}
		} else {
			if dbService.JobName.Valid && dbService.JobName.String != "" {
				serviceContext.JobName = dbService.JobName.String
			} else {
				serviceContext.JobName = service.Node.Name
			}

			serviceContext.Enabled = true
		}

		serviceContexts = append(serviceContexts, serviceContext)
	}

	response := ListServicesResponse{
		Services: serviceContexts,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshalling response: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

	return nil
}
