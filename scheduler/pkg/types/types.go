package types

import "time"

type Config struct {
	Port                      string        `env:"PORT" json:"port,omitempty"`
	DatabaseUrl               string        `env:"DATABASE_URL" json:"database_url,omitempty"`
	MessageBusUrl             string        `env:"MESSAGE_BUS_URL" json:"message_bus_url,omitempty"`
	CacheUrl                  string        `env:"CACHE_URL" json:"cache_url,omitempty"`
	WorkerUnackedMessageCount int           `env:"WORKER_UNACKED_MESSAGE_COUNT" json:"worker_unacked_message_count,omitempty"`
	WorkerStuckJobThreshold   time.Duration `env:"WORKER_STUCK_JOB_THRESHOLD" json:"worker_stuck_job_threshold,omitempty"`
	WorkerMaxJobRetries       int           `env:"WORKER_MAX_JOB_RETRIES" json:"worker_max_job_retries,omitempty"`
}
