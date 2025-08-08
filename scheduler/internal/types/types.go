package types

import "time"

type Config struct {
	Port                      string        `env:"PORT"`
	DatabaseUrl               string        `env:"DATABASE_URL"`
	MessageBusUrl             string        `env:"MESSAGE_BUS_URL"`
	CacheUrl                  string        `env:"CACHE_URL"`
	WorkerUnackedMessageCount int           `env:"WORKER_UNACKED_MESSAGE_COUNT"`
	WorkerStuckJobThreshold   time.Duration `env:"WORKER_STUCK_JOB_THRESHOLD"`
	WorkerMaxJobRetries       int           `env:"WORKER_MAX_JOB_RETRIES"`
}
