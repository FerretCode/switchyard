package types

import "time"

type Config struct {
	Port                      string        `env:"PORT" json:"port"`
	DatabaseUrl               string        `env:"DATABASE_URL" json:"database_url"`
	MessageBusUrl             string        `env:"MESSAGE_BUS_URL" json:"message_bus_url"`
	CacheUrl                  string        `env:"CACHE_URL" json:"cache_url"`
	WorkerUnackedMessageCount int           `env:"WORKER_UNACKED_MESSAGE_COUNT" json:"worker_unacked_message_count"`
	WorkerStuckJobThreshold   time.Duration `env:"WORKER_STUCK_JOB_THRESHOLD" json:"worker_stuck_job_threshold"`
	WorkerMaxJobRetries       int           `env:"WORKER_MAX_JOB_RETRIES" json:"worker_max_job_retries"`
}
