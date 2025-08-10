package configurator

type LocomotiveConfig struct {
	EnvironmentId    string   `env:"RAILWAY_ENVIRONMENT_ID"`
	IngestUrl        string   `env:"INGEST_URL"`
	RailwayApiKey    string   `env:"RAILWAY_API_KEY"`
	Train            []string `env:"TRAIN"`
	EnableDeployLogs bool     `env:"ENABLE_DEPLOY_LOGS"`
}
