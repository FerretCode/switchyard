package types

import "time"

type Config struct {
	Port                            string        `env:"PORT"`
	RailwayApiKey                   string        `env:"RAILWAY_API_KEY"`
	RailwayProjectId                string        `env:"RAILWAY_PROJECT_ID"`
	RailwayEnvironmentId            string        `env:"RAILWAY_ENVIRONMENT_ID"`
	RailwayMemoryUpscaleThreshold   float64       `env:"RAILWAY_MEMORY_UPSCALE_THRESHOLD"`
	RailwayCpuUpscaleThreshold      float64       `env:"RAILWAY_CPU_UPSCALE_THRESHOLD"`
	RailwayMemoryDownscaleThreshold float64       `env:"RAILWAY_MEMORY_DOWNSCALE_THRESHOLD"`
	RailwayCpuDownscaleThreshold    float64       `env:"RAILWAY_CPU_DOWNSCALE_THRESHOLD"`
	RailwaySelectedRegion           string        `env:"RAILWAY_SELECTED_REGION"`
	MonitoringInterval              time.Duration `env:"MONITORING_INTERVAL"`
	MetricHistorySize               int           `env:"METRIC_HISTORY_SIZE"`
	UpscaleCooldown                 time.Duration `env:"UPSCALE_COOLDOWN"`
	DownscaleCooldown               time.Duration `env:"DOWNSCALE_COOLDOWN"`
	MinReplicaCount                 int           `env:"MIN_REPLICA_COUNT"`
	MaxReplicaCount                 int           `env:"MAX_REPLICA_COUNT"`
	DatabaseUrl                     string        `env:"DATABASE_URL"`
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
