package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	messagebus "github.com/ferretcode/switchyard/incident/internal/message_bus"
	"github.com/ferretcode/switchyard/incident/internal/repositories"
	"github.com/ferretcode/switchyard/incident/pkg/types"
)

type WebhookService struct {
	Logger            *slog.Logger
	Config            *types.Config
	Queries           *repositories.Queries
	Context           context.Context
	MessageBusService *messagebus.MessageBusService
}

func NewWebhookService(logger *slog.Logger, config *types.Config, queries *repositories.Queries, context context.Context, messageBusService *messagebus.MessageBusService) WebhookService {
	return WebhookService{
		Logger:            logger,
		Config:            config,
		Queries:           queries,
		Context:           context,
		MessageBusService: messageBusService,
	}
}

func (w *WebhookService) SendGenericIncidentReport(message string) error {
	incidentReport := types.IncidentReport{
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	return w.sendIncidentReport(incidentReport)
}

func (w *WebhookService) SendDeploymentIncidentReport(message string, serviceId string, deploymentId string, environmentId string) error {
	incidentReport := types.IncidentReport{
		ServiceId:     serviceId,
		DeploymentId:  deploymentId,
		Message:       message,
		EnvironmentId: environmentId,
		Timestamp:     time.Now().Unix(),
	}

	return w.sendIncidentReport(incidentReport)
}

func (w *WebhookService) sendIncidentReport(incidentReport types.IncidentReport) error {
	_, err := w.Queries.CreateIncidentReport(w.Context, repositories.CreateIncidentReportParams{
		ServiceID:     incidentReport.ServiceId,
		DeploymentID:  incidentReport.DeploymentId,
		EnvironmentID: incidentReport.EnvironmentId,
		Message:       incidentReport.Message,
		Timestamp:     incidentReport.Timestamp,
	})
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}

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

	err = w.MessageBusService.SendIncidentReportMessage(incidentReport)
	if err != nil {
		return err
	}

	return nil
}
