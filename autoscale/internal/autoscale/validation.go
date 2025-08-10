package autoscale

import (
	"database/sql"
	"log/slog"

	"github.com/ferretcode/switchyard/autoscale/internal/repositories"
)

func populateUpsertServiceParams(upsertServiceRequest *UpsertServiceRequest) repositories.UpdateServiceParams {
	slog.Info("populating upsert service params", "request", upsertServiceRequest)

	params := repositories.UpdateServiceParams{
		ServiceID:                       *upsertServiceRequest.ServiceId,
		Enabled:                         defaultEnabled,
		RailwayMemoryUpscaleThreshold:   defaultMemoryUpscaleThreshold,
		RailwayCpuUpscaleThreshold:      defaultCpuUpscaleThreshold,
		RailwayMemoryDownscaleThreshold: defaultMemoryDownscaleThreshold,
		RailwayCpuDownscaleThreshold:    defaultCpuDownscaleThreshold,
		UpscaleCooldown:                 defaultUpscaleCooldown,
		DownscaleCooldown:               defaultDownscaleCooldown,
		MinReplicaCount:                 int32(defaultMinReplicaCount),
		MaxReplicaCount:                 int32(defaultMaxReplicaCount),
	}

	if upsertServiceRequest.JobName != nil {
		params.JobName = sql.NullString{String: *upsertServiceRequest.JobName, Valid: true}
	} else {
		params.JobName = sql.NullString{Valid: false}
	}

	if upsertServiceRequest.Enabled != nil {
		params.Enabled = *upsertServiceRequest.Enabled
	}

	if upsertServiceRequest.RailwayMemoryUpscaleThreshold != nil {
		params.RailwayMemoryUpscaleThreshold = *upsertServiceRequest.RailwayMemoryUpscaleThreshold
	}

	if upsertServiceRequest.RailwayCPUUpscaleThreshold != nil {
		params.RailwayCpuUpscaleThreshold = *upsertServiceRequest.RailwayCPUUpscaleThreshold
	}

	if upsertServiceRequest.RailwayMemoryDownscaleThreshold != nil {
		params.RailwayMemoryDownscaleThreshold = *upsertServiceRequest.RailwayMemoryDownscaleThreshold
	}

	if upsertServiceRequest.RailwayCPUDownscaleThreshold != nil {
		params.RailwayCpuDownscaleThreshold = *upsertServiceRequest.RailwayCPUDownscaleThreshold
	}

	if upsertServiceRequest.UpscaleCooldown != nil {
		params.UpscaleCooldown = *upsertServiceRequest.UpscaleCooldown
	}

	if upsertServiceRequest.DownscaleCooldown != nil {
		params.DownscaleCooldown = *upsertServiceRequest.DownscaleCooldown
	}

	if upsertServiceRequest.MinReplicaCount != nil {
		params.MinReplicaCount = int32(*upsertServiceRequest.MinReplicaCount)
	}

	if upsertServiceRequest.MaxReplicaCount != nil {
		params.MaxReplicaCount = int32(*upsertServiceRequest.MaxReplicaCount)
	}

	return params
}

func populateRegisterServiceParams(req *RegisterServiceRequest) repositories.CreateServiceParams {
	params := repositories.CreateServiceParams{
		ServiceID:                       *req.ServiceId,
		Enabled:                         defaultEnabled,
		RailwayMemoryUpscaleThreshold:   defaultMemoryUpscaleThreshold,
		RailwayCpuUpscaleThreshold:      defaultCpuUpscaleThreshold,
		RailwayMemoryDownscaleThreshold: defaultMemoryDownscaleThreshold,
		RailwayCpuDownscaleThreshold:    defaultCpuDownscaleThreshold,
		UpscaleCooldown:                 defaultUpscaleCooldown,
		DownscaleCooldown:               defaultDownscaleCooldown,
		MinReplicaCount:                 int32(defaultMinReplicaCount),
		MaxReplicaCount:                 int32(defaultMaxReplicaCount),
	}

	if req.JobName != nil {
		params.JobName = sql.NullString{String: *req.JobName, Valid: true}
	} else {
		params.JobName = sql.NullString{Valid: false}
	}

	if req.Enabled != nil {
		params.Enabled = *req.Enabled
	}

	if req.RailwayMemoryUpscaleThreshold != nil {
		params.RailwayMemoryUpscaleThreshold = *req.RailwayMemoryUpscaleThreshold
	}

	if req.RailwayCPUUpscaleThreshold != nil {
		params.RailwayCpuUpscaleThreshold = *req.RailwayCPUUpscaleThreshold
	}

	if req.RailwayMemoryDownscaleThreshold != nil {
		params.RailwayMemoryDownscaleThreshold = *req.RailwayMemoryDownscaleThreshold
	}

	if req.RailwayCPUDownscaleThreshold != nil {
		params.RailwayCpuDownscaleThreshold = *req.RailwayCPUDownscaleThreshold
	}

	if req.UpscaleCooldown != nil {
		params.UpscaleCooldown = *req.UpscaleCooldown
	}

	if req.DownscaleCooldown != nil {
		params.DownscaleCooldown = *req.DownscaleCooldown
	}

	if req.MinReplicaCount != nil {
		params.MinReplicaCount = int32(*req.MinReplicaCount)
	}

	if req.MaxReplicaCount != nil {
		params.MaxReplicaCount = int32(*req.MaxReplicaCount)
	}

	return params
}

func applyDefaultsForCreate(params repositories.CreateServiceParams) repositories.CreateServiceParams {
	if !params.Enabled {
		params.Enabled = true
	}
	if params.RailwayMemoryUpscaleThreshold == 0 {
		params.RailwayMemoryUpscaleThreshold = defaultMemoryUpscaleThreshold
	}
	if params.RailwayCpuUpscaleThreshold == 0 {
		params.RailwayCpuUpscaleThreshold = defaultCpuUpscaleThreshold
	}
	if params.RailwayMemoryDownscaleThreshold == 0 {
		params.RailwayMemoryDownscaleThreshold = defaultMemoryDownscaleThreshold
	}
	if params.RailwayCpuDownscaleThreshold == 0 {
		params.RailwayCpuDownscaleThreshold = defaultCpuDownscaleThreshold
	}
	if params.UpscaleCooldown == "" {
		params.UpscaleCooldown = defaultUpscaleCooldown
	}
	if params.DownscaleCooldown == "" {
		params.DownscaleCooldown = defaultDownscaleCooldown
	}
	if params.MinReplicaCount == 0 {
		params.MinReplicaCount = defaultMinReplicaCount
	}
	if params.MaxReplicaCount == 0 {
		params.MaxReplicaCount = defaultMaxReplicaCount
	}
	return params
}
