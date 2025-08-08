package types

import "time"

type Config struct {
	Port                            string        `env:"PORT" json:"port"`
	RailwayApiKey                   string        `env:"RAILWAY_API_KEY" json:"railway_api_key"`
	RailwayProjectId                string        `env:"RAILWAY_PROJECT_ID" json:"railway_project_id"`
	RailwayEnvironmentId            string        `env:"RAILWAY_ENVIRONMENT_ID" json:"railway_environment_id"`
	RailwayMemoryUpscaleThreshold   float64       `env:"RAILWAY_MEMORY_UPSCALE_THRESHOLD" json:"railway_memory_upscale_threshold"`
	RailwayCpuUpscaleThreshold      float64       `env:"RAILWAY_CPU_UPSCALE_THRESHOLD" json:"railway_cpu_upscale_threshold"`
	RailwayMemoryDownscaleThreshold float64       `env:"RAILWAY_MEMORY_DOWNSCALE_THRESHOLD" json:"railway_memory_downscale_threshold"`
	RailwayCpuDownscaleThreshold    float64       `env:"RAILWAY_CPU_DOWNSCALE_THRESHOLD" json:"railway_cpu_downscale_threshold"`
	RailwaySelectedRegion           string        `env:"RAILWAY_SELECTED_REGION" json:"railway_selected_region"`
	MonitoringInterval              time.Duration `env:"MONITORING_INTERVAL" json:"monitoring_interval"`
	MetricHistorySize               int           `env:"METRIC_HISTORY_SIZE" json:"metric_history_size"`
	UpscaleCooldown                 time.Duration `env:"UPSCALE_COOLDOWN" json:"upscale_cooldown"`
	DownscaleCooldown               time.Duration `env:"DOWNSCALE_COOLDOWN" json:"downscale_cooldown"`
	MinReplicaCount                 int           `env:"MIN_REPLICA_COUNT" json:"min_replica_count"`
	MaxReplicaCount                 int           `env:"MAX_REPLICA_COUNT" json:"max_replica_count"`
	DatabaseUrl                     string        `env:"DATABASE_URL" json:"database_url"`
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
