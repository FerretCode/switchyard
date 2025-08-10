package ingest

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ferretcode/switchyard/incident/internal/webhook"
	"github.com/ferretcode/switchyard/incident/pkg/types"
)

var lastIncidentSent time.Time

type IngestService struct {
	Logger             *slog.Logger
	Config             *types.Config
	IncidentStats      *types.IncidentStats
	PrometheusCounters *types.PrometheusCounters
	WebhookService     *webhook.WebhookService
}

func NewIngestService(logger *slog.Logger, incidentStats *types.IncidentStats, config *types.Config, prometheusCounters *types.PrometheusCounters, webhookService *webhook.WebhookService) IngestService {
	return IngestService{
		Logger:             logger,
		IncidentStats:      incidentStats,
		Config:             config,
		PrometheusCounters: prometheusCounters,
		WebhookService:     webhookService,
	}
}

func (i *IngestService) Ingest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var loglines []types.Logline

	if err := json.Unmarshal(bytes, &loglines); err != nil {
		return err
	}

	for _, log := range loglines {
		i.PrometheusCounters.LoglineCounter.Inc()

		if strings.ToUpper(log.Severity) != "ERROR" {
			continue
		}

		i.PrometheusCounters.ErrorCounter.Inc()
		i.IncidentStats.ErrorCount++

		if !i.detectIncident() {
			continue
		}

		i.Logger.Info("incident detected")

		now := time.Now()
		cutoff := now.Add(-i.Config.IncidentAnalysisWindow)

		if lastIncidentSent.Before(cutoff) {
			if log.Metadata["_service_id"] != "" && log.Metadata["_deployment_id"] != "" && log.Metadata["_environment_id"] != "" {
				err := i.WebhookService.SendDeploymentIncidentReport(log.Message, log.Metadata["_service_id"], log.Metadata["_deployment_id"], log.Metadata["_environment_id"])
				if err != nil {
					return err
				}
			} else {
				err := i.WebhookService.SendGenericIncidentReport(log.Message)
				if err != nil {
					return err
				}
			}
			lastIncidentSent = now
		}
	}

	w.WriteHeader(200)

	return nil
}

func (i *IngestService) detectIncident() bool {
	now := time.Now()
	window := i.Config.IncidentAnalysisWindow
	cutoff := now.Add(-window)

	i.IncidentStats.ErrorWindow = filterRecent(i.IncidentStats.ErrorWindow, cutoff)
	i.IncidentStats.ErrorWindow = append(i.IncidentStats.ErrorWindow, now)

	if len(i.IncidentStats.ErrorWindow) > i.Config.IncidentAnalysisErrorThreshold {
		return true
	}

	return false
}

func filterRecent(timestamps []time.Time, cutoff time.Time) []time.Time {
	var recent []time.Time
	for _, timestamp := range timestamps {
		if timestamp.After(cutoff) {
			recent = append(recent, timestamp)
		}
	}

	return recent
}
