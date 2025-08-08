package messagebus

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ferretcode/switchyard/scheduler/pkg/types"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type MessageBusService struct {
	Logger    *slog.Logger
	Config    *types.Config
	Conn      *amqp.Connection
	RedisConn *redis.Client
	Context   context.Context
}

func NewMessageBusService(logger *slog.Logger, conn *amqp.Connection, config *types.Config, redisConn *redis.Client, context context.Context) MessageBusService {
	return MessageBusService{
		Logger:    logger,
		Conn:      conn,
		RedisConn: redisConn,
		Config:    config,
		Context:   context,
	}
}

func (m *MessageBusService) SendRetryJobMessage(jobName string, jobContext map[string]any, jobId string) error {
	channel, queue, err := m.declareQueueAndChannel()
	if err != nil {
		return err
	}
	defer channel.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	scheduleJobBody := map[string]any{
		"job_name":    jobName,
		"job_context": jobContext,
		"job_id":      jobId,
	}

	bodyBytes, err := json.Marshal(scheduleJobBody)
	if err != nil {
		return err
	}

	m.Logger.Info("publishing job receipt to message queue", "job-id", jobId)

	err = channel.PublishWithContext(ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         bodyBytes,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		return err
	}

	return nil
}

func (m *MessageBusService) SendScheduleJobMessage(jobName string, jobContext map[string]any) error {
	channel, queue, err := m.declareQueueAndChannel()
	if err != nil {
		return err
	}
	defer channel.Close()

	contextBytes, err := json.Marshal(jobContext)
	if err != nil {
		return err
	}

	jobId := uuid.NewString()

	err = m.RedisConn.HSet(m.Context, "jobs:"+jobId, map[string]interface{}{
		"status":      "pending",
		"created_at":  time.Now().Unix(),
		"updated_at":  time.Now().Unix(),
		"retry_count": 0,
		"message":     "",
		"job_name":    jobName,
		"job_context": string(contextBytes),
	}).Err()
	if err != nil {
		return err
	}

	err = m.RedisConn.ZAdd(m.Context, "jobs:pending", redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: jobId,
	}).Err()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	scheduleJobBody := map[string]any{
		"job_name":    jobName,
		"job_context": jobContext,
		"job_id":      jobId,
	}

	bodyBytes, err := json.Marshal(scheduleJobBody)
	if err != nil {
		return err
	}

	m.Logger.Info("publishing job receipt to message queue", "job-id", jobId)

	err = channel.PublishWithContext(ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         bodyBytes,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		return err
	}

	return nil
}

func (m *MessageBusService) SubscribeToJobFinishedMessages() (chan bool, error) {
	channel, err := m.Conn.Channel()
	if err != nil {
		return nil, err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"jobs-finished",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	done := make(chan bool)

	for {
		select {
		case <-done:
			return done, nil
		case delivery := <-msgs:
			m.handleMessageFinishedDelivery(delivery)
		}
	}
}

func (m *MessageBusService) handleMessageFinishedDelivery(delivery amqp.Delivery) {
	finishJobMessage := FinishJobMessage{}

	if err := json.Unmarshal(delivery.Body, &finishJobMessage); err != nil {
		m.Logger.Error("error decoding message", "err", err)
		return
	}

	jobKey := "jobs:" + finishJobMessage.JobId

	updatedStatus := ""

	switch finishJobMessage.Status {
	case OK:
		updatedStatus = "ok"
	case ERROR:
		updatedStatus = "error"
	}

	err := m.RedisConn.HSet(m.Context, jobKey, "status", updatedStatus).Err()
	if err != nil {
		m.Logger.Error("error updating job status", "err", err)
		return
	}

	err = m.RedisConn.HSet(m.Context, jobKey, "message", finishJobMessage.Message).Err()
	if err != nil {
		m.Logger.Error("error updating job message", "err", err)
		return
	}

	m.Logger.Info("job has been processed successfully", "job-id", finishJobMessage.JobId)

	delivery.Ack(false)
}

func (m *MessageBusService) declareQueueAndChannel() (*amqp.Channel, amqp.Queue, error) {
	channel, err := m.Conn.Channel()
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, err
	}

	channel.Qos(
		m.Config.WorkerUnackedMessageCount,
		0,
		false,
	)

	queue, err := channel.QueueDeclare(
		"jobs",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, err
	}

	return channel, queue, nil
}
