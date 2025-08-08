package types

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Port                                  string        `env:"PORT"`
	IncidentAnalysisWindow                time.Duration `env:"INCIDENT_ANALYSIS_WINDOW"`
	IncidentAnalysisErrorThreshold        int           `env:"INCIDENT_ANALYSIS_ERROR_THRESHOLD"`
	ServiceMonitorPollingRate             int           `env:"SERVICE_MONITOR_POLLING_RATE"`
	ServiceMonitorPollingTimeout          int           `env:"SERVICE_MONITOR_POLLING_TIMEOUT"`
	ServiceMonitorInterestedStatusChanges []string      `env:"SERVICE_MONITOR_INTERESTED_STATUS_CHANGES"`
	RailwayApiKey                         string        `env:"RAILWAY_API_KEY"`
	RailwayEnvironmentId                  string        `env:"RAILWAY_ENVIRONMENT_ID"`
	RailwayServiceIds                     []string      `env:"RAILWAY_SERVICE_IDS"`
	IncidentReportWebhookUrl              string        `env:"INCIDENT_REPORT_WEBHOOK_URL"`
	IncidentReportAdditionalHeaders       string        `env:"INCIDENT_REPORT_ADDITIONAL_HEADERS"`
}

type Logline struct {
	Message   string `json:"message"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
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
	ServiceId    string `json:"service_id"`
	DeploymentId string `json:"deployment_id"`
	Message      string `json:"message"`
	Timestamp    int64  `json:"timestamp"`
}
