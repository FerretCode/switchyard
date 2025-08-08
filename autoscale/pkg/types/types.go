package types

import "time"

type Config struct {
	Port                            string        `env:"PORT" json:"port,omitempty"`
	RailwayApiKey                   string        `env:"RAILWAY_API_KEY" json:"railway_api_key,omitempty"`
	RailwayProjectId                string        `env:"RAILWAY_PROJECT_ID" json:"railway_project_id,omitempty"`
	RailwayEnvironmentId            string        `env:"RAILWAY_ENVIRONMENT_ID" json:"railway_environment_id,omitempty"`
	RailwayMemoryUpscaleThreshold   float64       `env:"RAILWAY_MEMORY_UPSCALE_THRESHOLD" json:"railway_memory_upscale_threshold,omitempty"`
	RailwayCpuUpscaleThreshold      float64       `env:"RAILWAY_CPU_UPSCALE_THRESHOLD" json:"railway_cpu_upscale_threshold,omitempty"`
	RailwayMemoryDownscaleThreshold float64       `env:"RAILWAY_MEMORY_DOWNSCALE_THRESHOLD" json:"railway_memory_downscale_threshold,omitempty"`
	RailwayCpuDownscaleThreshold    float64       `env:"RAILWAY_CPU_DOWNSCALE_THRESHOLD" json:"railway_cpu_downscale_threshold,omitempty"`
	RailwaySelectedRegion           string        `env:"RAILWAY_SELECTED_REGION" json:"railway_selected_region,omitempty"`
	MonitoringInterval              time.Duration `env:"MONITORING_INTERVAL" json:"monitoring_interval,omitempty"`
	MetricHistorySize               int           `env:"METRIC_HISTORY_SIZE" json:"metric_history_size,omitempty"`
	UpscaleCooldown                 time.Duration `env:"UPSCALE_COOLDOWN" json:"upscale_cooldown,omitempty"`
	DownscaleCooldown               time.Duration `env:"DOWNSCALE_COOLDOWN" json:"downscale_cooldown,omitempty"`
	MinReplicaCount                 int           `env:"MIN_REPLICA_COUNT" json:"min_replica_count,omitempty"`
	MaxReplicaCount                 int           `env:"MAX_REPLICA_COUNT" json:"max_replica_count,omitempty"`
	DatabaseUrl                     string        `env:"DATABASE_URL" json:"database_url,omitempty"`
}

type MetricHistory struct {
	CPU    []float64
	Memory []float64
	Times  []time.Time
}

type ServiceHistoryCache struct {
	ServiceHistories map[string]MetricHistory
}

type ScalingContext struct {
	CpuPercent      float64
	MemPercent      float64
	AvgCpu          float64
	AvgMem          float64
	HasCpuSpike     bool
	HasMemSpike     bool
	CpuTrend        float64
	MemTrend        float64
	CurrentReplicas int
	Now             time.Time
}
