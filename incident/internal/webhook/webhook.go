package webhook

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ferretcode/switchyard/incident/internal/types"
)

type WebhookService struct {
	Logger *slog.Logger
	Config *types.Config
}

func NewWebhookService(logger *slog.Logger, config *types.Config) WebhookService {
	return WebhookService{
		Logger: logger,
		Config: config,
	}
}

func (w *WebhookService) SendGenericIncidentReport(message string) error {
	incidentReport := types.IncidentReport{
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	return w.sendIncidentReport(incidentReport)
}

func (w *WebhookService) SendDeploymentIncidentReport(message string, serviceId string, deploymentId string) error {
	incidentReport := types.IncidentReport{
		ServiceId:    serviceId,
		DeploymentId: deploymentId,
		Message:      message,
		Timestamp:    time.Now().Unix(),
	}

	return w.sendIncidentReport(incidentReport)
}

func (w *WebhookService) sendIncidentReport(incidentReport types.IncidentReport) error {
	incidentReportBytes, err := json.Marshal(incidentReport)
	if err != nil {
		return err
	}

	w.Logger.Info("sending incident report to url", "url", w.Config.IncidentReportWebhookUrl)

	req, err := http.NewRequest(
		"POST",
		w.Config.IncidentReportWebhookUrl,
		bytes.NewReader(incidentReportBytes),
	)

	headers := strings.SplitSeq(w.Config.IncidentReportAdditionalHeaders, ";")
	for header := range headers {
		pair := strings.Split(header, "=")

		if len(pair) == 2 {
			req.Header.Set(pair[0], pair[1])
		}
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
