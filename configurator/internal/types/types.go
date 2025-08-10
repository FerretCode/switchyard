package types

type Config struct {
	Port                 string `env:"PORT"`
	Environment          string `env:"ENVIRONMENT"`
	DatabaseUrl          string `env:"DATABASE_URL"`
	RailwayApiKey        string `env:"RAILWAY_API_KEY"`
	RailwayProjectId     string `env:"RAILWAY_PROJECT_ID"`
	RailwayEnvironmentId string `env:"RAILWAY_ENVIRONMENT_ID"`
}
