package types

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Port                                  string        `env:"PORT" json:"port,omitempty"`
	DatabaseUrl                           string        `env:"DATABASE_URL" json:"database_url,omitempty"`
	MessageBusUrl                         string        `env:"MESSAGE_BUS_URL" json:"message_bus_url,omitempty"`
	IncidentAnalysisWindow                time.Duration `env:"INCIDENT_ANALYSIS_WINDOW" json:"incident_analysis_window,omitempty"`
	IncidentAnalysisErrorThreshold        int           `env:"INCIDENT_ANALYSIS_ERROR_THRESHOLD" json:"incident_analysis_error_threshold,omitempty"`
	ServiceMonitorPollingRate             time.Duration `env:"SERVICE_MONITOR_POLLING_RATE" json:"service_monitor_polling_rate,omitempty"`
	ServiceMonitorPollingTimeout          time.Duration `env:"SERVICE_MONITOR_POLLING_TIMEOUT" json:"service_monitor_polling_timeout,omitempty"`
	ServiceMonitorInterestedStatusChanges []string      `env:"SERVICE_MONITOR_INTERESTED_STATUS_CHANGES" json:"service_monitor_interested_status_changes,omitempty"`
	RailwayApiKey                         string        `env:"RAILWAY_API_KEY" json:"railway_api_key,omitempty"`
	RailwayEnvironmentId                  string        `env:"RAILWAY_ENVIRONMENT_ID" json:"railway_environment_id,omitempty"`
	RailwayServiceIds                     []string      `env:"RAILWAY_SERVICE_IDS" json:"railway_service_ids,omitempty"`
	IncidentReportWebhookUrl              string        `env:"INCIDENT_REPORT_WEBHOOK_URL" json:"incident_report_webhook_url,omitempty"`
	IncidentReportAdditionalHeaders       string        `env:"INCIDENT_REPORT_ADDITIONAL_HEADERS" json:"incident_report_additional_headers,omitempty"`
}

type Logline struct {
	Message   string            `json:"message"`
	Severity  string            `json:"severity"`
	Timestamp string            `json:"timestamp"`
	Metadata  map[string]string `json:"_metadata,omitempty"`
}

type IncidentStats struct {
	ErrorCount  int
	LastError   time.Time
	ErrorWindow []time.Time
}

type PrometheusCounters struct {
	LoglineCounter prometheus.Counter
	ErrorCounter   prometheus.Counter
}

type DeploymentCache struct {
	Deployments map[string]string
	Mutex       sync.Mutex
}

type IncidentReport struct {
	ServiceId     string `json:"service_id"`
	DeploymentId  string `json:"deployment_id"`
	EnvironmentId string `json:"environment_id"`
	Message       string `json:"message"`
	Timestamp     int64  `json:"timestamp"`
}
