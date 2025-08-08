package types

type Config struct {
	Port               string `env:"PORT"`
	SessionsCookieName string `env:"SESSIONS_COOKIE_NAME"`
	AdminUsername      string `env:"ADMIN_USERNAME"`
	AdminPassword      string `env:"ADMIN_PASSWORD"`

	AutoscaleServiceUrl         string `env:"AUTOSCALE_SERVICE_URL"`
	FeatureFlagsServiceUrl      string `env:"FEATURE_FLAGS_SERVICE_URL"`
	IncidentReportingServiceUrl string `env:"INCIDENT_REPORTING_SERVICE_URL"`
	SchedulerServiceUrl         string `env:"SCHEDULER_SERVICE_URL"`
}
