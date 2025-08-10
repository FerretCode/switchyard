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

func (a *AutoscaleService) UpsertService(w http.ResponseWriter, r *http.Request) error {
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	var upsertServiceRequest UpsertServiceRequest
	if err := json.Unmarshal(requestBytes, &upsertServiceRequest); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}

	if upsertServiceRequest.ServiceId == nil || *upsertServiceRequest.ServiceId == "" {
		http.Error(w, "Service ID is required", http.StatusBadRequest)
		return nil
	}

	_, err = a.Queries.GetService(a.Context, *upsertServiceRequest.ServiceId)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error fetching service from database: %w", err)
	}

	params := populateUpsertServiceParams(&upsertServiceRequest)

	if err == sql.ErrNoRows {
		createParams := repositories.CreateServiceParams{
			ServiceID: *upsertServiceRequest.ServiceId,
			JobName:   sql.NullString{Valid: false},
			Enabled:   false,
		}

		if upsertServiceRequest.JobName != nil {
			createParams.JobName = sql.NullString{String: *upsertServiceRequest.JobName, Valid: true}
		}
		if upsertServiceRequest.Enabled != nil {
			createParams.Enabled = *upsertServiceRequest.Enabled
		}
		if upsertServiceRequest.RailwayMemoryUpscaleThreshold != nil {
			createParams.RailwayMemoryUpscaleThreshold = *upsertServiceRequest.RailwayMemoryUpscaleThreshold
		}
		if upsertServiceRequest.RailwayCPUUpscaleThreshold != nil {
			createParams.RailwayCpuUpscaleThreshold = *upsertServiceRequest.RailwayCPUUpscaleThreshold
		}
		if upsertServiceRequest.RailwayMemoryDownscaleThreshold != nil {
			createParams.RailwayMemoryDownscaleThreshold = *upsertServiceRequest.RailwayMemoryDownscaleThreshold
		}
		if upsertServiceRequest.RailwayCPUDownscaleThreshold != nil {
			createParams.RailwayCpuDownscaleThreshold = *upsertServiceRequest.RailwayCPUDownscaleThreshold
		}
		if upsertServiceRequest.UpscaleCooldown != nil {
			createParams.UpscaleCooldown = *upsertServiceRequest.UpscaleCooldown
		}
		if upsertServiceRequest.DownscaleCooldown != nil {
			createParams.DownscaleCooldown = *upsertServiceRequest.DownscaleCooldown
		}
		if upsertServiceRequest.MinReplicaCount != nil {
			createParams.MinReplicaCount = int32(*upsertServiceRequest.MinReplicaCount)
		}
		if upsertServiceRequest.MaxReplicaCount != nil {
			createParams.MaxReplicaCount = int32(*upsertServiceRequest.MaxReplicaCount)
		}

		createParams = applyDefaultsForCreate(createParams)

		_, err = a.Queries.CreateService(a.Context, createParams)
		if err != nil {
			return fmt.Errorf("error creating service: %w", err)
		}

		w.WriteHeader(http.StatusOK)
		return nil
	}

	_, err = a.Queries.UpdateService(a.Context, params)
	if err != nil {
		return fmt.Errorf("error updating service: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (a *AutoscaleService) RegisterService(w http.ResponseWriter, r *http.Request) error {
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	var registerServiceRequest RegisterServiceRequest
	if err := json.Unmarshal(requestBytes, &registerServiceRequest); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}

	if registerServiceRequest.ServiceId == nil || *registerServiceRequest.ServiceId == "" {
		http.Error(w, "Service ID is required", http.StatusBadRequest)
		return nil
	}

	_, err = a.Queries.GetService(a.Context, *registerServiceRequest.ServiceId)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error fetching service from database: %w", err)
	}

	if err == nil {
		a.Logger.Info("service already exists, updating instead of registering", "service_id", registerServiceRequest.ServiceId)
		return a.UpsertService(w, r)
	}

	params := populateRegisterServiceParams(&registerServiceRequest)
	params = applyDefaultsForCreate(params)

	_, err = a.Queries.CreateService(a.Context, params)
	if err != nil {
		return fmt.Errorf("error creating service: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (a *AutoscaleService) UnregisterService(w http.ResponseWriter, r *http.Request) error {
	serviceId := chi.URLParam(r, "id")

	existingService, err := a.Queries.GetService(a.Context, serviceId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("service with id %s does not exist", serviceId)
		}
		return fmt.Errorf("error fetching service from database: %w", err)
	}

	// we want to keep the already existing service options, so we'll just disable it
	if existingService.Enabled {
		_, err = a.Queries.SetServiceEnabled(a.Context, repositories.SetServiceEnabledParams{
			ServiceID: serviceId,
			Enabled:   false,
		})
		if err != nil {
			return fmt.Errorf("error disabling service: %w", err)
		}

		w.WriteHeader(200)
		return nil
	}

	err = a.Queries.DeleteService(a.Context, serviceId)
	if err != nil {
		return fmt.Errorf("error deleting service: %w", err)
	}

	w.WriteHeader(200)
	return nil
}

func (a *AutoscaleService) SetServiceEnabled(w http.ResponseWriter, r *http.Request) error {
	serviceId := chi.URLParam(r, "id")

	existingService, err := a.Queries.GetService(a.Context, serviceId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("service with id %s does not exist", serviceId)
		}
		return fmt.Errorf("error fetching service from database: %w", err)
	}

	enabled := !existingService.Enabled

	_, err = a.Queries.SetServiceEnabled(a.Context, repositories.SetServiceEnabledParams{
		ServiceID: serviceId,
		Enabled:   enabled,
	})
	if err != nil {
		return fmt.Errorf("error setting service enabled state: %w", err)
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

		serviceContext.ServiceName = service.Node.Name
		serviceContext.ServiceId = service.Node.Id
		serviceContext.ProjectId = project.Project.Id
		serviceContext.EnvironmentId = a.Config.RailwayEnvironmentId
		serviceContext.EnvironmentName = currentEnvironmentName
		serviceContext.Replicas = currentReplicas
		serviceContext.LastScaledAt = service.Node.ServiceInstances.Edges[0].Node.LatestDeployment.CreatedAt

		dbService, err := a.Queries.GetService(a.Context, service.Node.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				serviceContext.Enabled = false

				serviceContext.MinReplicas = a.Config.MinReplicaCount
				serviceContext.MaxReplicas = a.Config.MaxReplicaCount
				serviceContext.CpuUpscaleThreshold = a.Config.RailwayCpuUpscaleThreshold * 100
				serviceContext.MemoryUpscaleThreshold = a.Config.RailwayMemoryUpscaleThreshold * 100
				serviceContext.CpuDownscaleThreshold = a.Config.RailwayCpuDownscaleThreshold * 100
				serviceContext.MemoryDownscaleThreshold = a.Config.RailwayMemoryDownscaleThreshold * 100
				serviceContext.UpscaleCooldown = a.Config.UpscaleCooldown.String()
				serviceContext.DownscaleCooldown = a.Config.DownscaleCooldown.String()
			} else {
				return fmt.Errorf("error fetching service from database: %w", err)
			}
		} else {
			if dbService.JobName.Valid && dbService.JobName.String != "" {
				serviceContext.JobName = dbService.JobName.String
			}

			serviceContext.MinReplicas = int(dbService.MinReplicaCount)
			serviceContext.MaxReplicas = int(dbService.MaxReplicaCount)
			serviceContext.CpuUpscaleThreshold = dbService.RailwayCpuUpscaleThreshold * 100
			serviceContext.MemoryUpscaleThreshold = dbService.RailwayMemoryUpscaleThreshold * 100
			serviceContext.CpuDownscaleThreshold = dbService.RailwayCpuDownscaleThreshold * 100
			serviceContext.MemoryDownscaleThreshold = dbService.RailwayMemoryDownscaleThreshold * 100
			serviceContext.UpscaleCooldown = dbService.UpscaleCooldown
			serviceContext.DownscaleCooldown = dbService.DownscaleCooldown

			serviceContext.Enabled = dbService.Enabled
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
