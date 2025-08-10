package autoscale

import (
	"github.com/ferretcode/switchyard/autoscale/internal/repositories"
)

const (
	defaultEnabled                  = true
	defaultMemoryUpscaleThreshold   = 0.80
	defaultCpuUpscaleThreshold      = 0.80
	defaultMemoryDownscaleThreshold = 0.20
	defaultCpuDownscaleThreshold    = 0.20
	defaultUpscaleCooldown          = "1m"
	defaultDownscaleCooldown        = "2m"
	defaultMinReplicaCount          = 1
	defaultMaxReplicaCount          = 10
)

type UpsertServiceRequest struct {
	ServiceId                       *string  `json:"service_id,omitempty"`
	JobName                         *string  `json:"job_name,omitempty"`
	Enabled                         *bool    `json:"enabled,omitempty"`
	RailwayMemoryUpscaleThreshold   *float64 `json:"railway_memory_upscale_threshold,omitempty"`
	RailwayCPUUpscaleThreshold      *float64 `json:"railway_cpu_upscale_threshold,omitempty"`
	RailwayMemoryDownscaleThreshold *float64 `json:"railway_memory_downscale_threshold,omitempty"`
	RailwayCPUDownscaleThreshold    *float64 `json:"railway_cpu_downscale_threshold,omitempty"`
	UpscaleCooldown                 *string  `json:"upscale_cooldown,omitempty"`
	DownscaleCooldown               *string  `json:"downscale_cooldown,omitempty"`
	MinReplicaCount                 *int     `json:"min_replica_count,omitempty"`
	MaxReplicaCount                 *int     `json:"max_replica_count,omitempty"`
}

type RegisterServiceRequest struct {
	ServiceId                       *string  `json:"service_id,omitempty"`
	JobName                         *string  `json:"job_name,omitempty"`
	Enabled                         *bool    `json:"enabled,omitempty"`
	RailwayMemoryUpscaleThreshold   *float64 `json:"railway_memory_upscale_threshold,omitempty"`
	RailwayCPUUpscaleThreshold      *float64 `json:"railway_cpu_upscale_threshold,omitempty"`
	RailwayMemoryDownscaleThreshold *float64 `json:"railway_memory_downscale_threshold,omitempty"`
	RailwayCPUDownscaleThreshold    *float64 `json:"railway_cpu_downscale_threshold,omitempty"`
	UpscaleCooldown                 *string  `json:"upscale_cooldown,omitempty"`
	DownscaleCooldown               *string  `json:"downscale_cooldown,omitempty"`
	MinReplicaCount                 *int     `json:"min_replica_count,omitempty"`
	MaxReplicaCount                 *int     `json:"max_replica_count,omitempty"`
}

type ValidService struct {
	ServiceId string               `json:"service_id"`
	Service   repositories.Service `json:"service"`
}

type ListServicesResponse struct {
	Services []ServiceContext `json:"services"`
}

type ServiceContext struct {
	ServiceId                string  `json:"service_id"`
	ProjectId                string  `json:"project_id"`
	JobName                  string  `json:"job_name"`
	ServiceName              string  `json:"service_name"`
	EnvironmentName          string  `json:"environment_name"`
	EnvironmentId            string  `json:"environment_id"`
	Replicas                 int     `json:"replicas"`
	MinReplicas              int     `json:"min_replicas"`
	MaxReplicas              int     `json:"max_replicas"`
	CpuUpscaleThreshold      float64 `json:"cpu_upscale_threshold"`
	MemoryUpscaleThreshold   float64 `json:"memory_upscale_threshold"`
	CpuDownscaleThreshold    float64 `json:"cpu_downscale_threshold"`
	MemoryDownscaleThreshold float64 `json:"memory_downscale_threshold"`
	UpscaleCooldown          string  `json:"upscale_cooldown"`
	DownscaleCooldown        string  `json:"downscale_cooldown"`
	LastScaledAt             string  `json:"last_scaled_at"`
	Enabled                  bool    `json:"enabled"`
}
