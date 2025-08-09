package messagebus

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ferretcode/switchyard/incident/pkg/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBusService struct {
	Logger  *slog.Logger
	Config  *types.Config
	Conn    *amqp.Connection
	Context context.Context
}

func NewMessageBusService(logger *slog.Logger, conn *amqp.Connection, config *types.Config, context context.Context) MessageBusService {
	return MessageBusService{
		Logger:  logger,
		Conn:    conn,
		Config:  config,
		Context: context,
	}
}

func (m *MessageBusService) SendIncidentReportMessage(incidentReport types.IncidentReport) error {
	channel, queue, err := m.declareQueueAndChannel()
	if err != nil {
		return err
	}
	defer channel.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bodyBytes, err := json.Marshal(incidentReport)
	if err != nil {
		return err
	}

	m.Logger.Info("publishing incident report to message queue", "incident-report", incidentReport)

	err = channel.PublishWithContext(ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MessageBusService) declareQueueAndChannel() (*amqp.Channel, amqp.Queue, error) {
	channel, err := m.Conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	queue, err := channel.QueueDeclare(
		"incident-reports",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return channel, queue, nil
}
