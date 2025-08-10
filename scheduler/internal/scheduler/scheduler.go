package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	messagebus "github.com/ferretcode/switchyard/scheduler/internal/message_bus"
	"github.com/ferretcode/switchyard/scheduler/internal/repositories"
	"github.com/go-chi/chi/v5"
)

type SchedulerService struct {
	Logger            *slog.Logger
	Queries           *repositories.Queries
	Context           context.Context
	MessageBusService *messagebus.MessageBusService
}

func NewSchedulerService(logger *slog.Logger, queries *repositories.Queries, context context.Context, messageBusService *messagebus.MessageBusService) SchedulerService {
	return SchedulerService{
		Logger:            logger,
		Queries:           queries,
		Context:           context,
		MessageBusService: messageBusService,
	}
}

func (s *SchedulerService) GetJobStatistics(w http.ResponseWriter, r *http.Request) error {
	jobName := chi.URLParam(r, "name")

	jobStatistics, err := s.Queries.AggregateJobReceiptsByJobID(s.Context, jobName)
	if err != nil {
		return fmt.Errorf("error fetching job statistics: %w", err)
	}

	responseBytes, err := json.Marshal(jobStatistics)
	if err != nil {
		return fmt.Errorf("error encoding job statistics: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

	return nil
}

func (s *SchedulerService) RegisterWorkerService(w http.ResponseWriter, r *http.Request) error {
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	registerWorkerServiceRequest := RegisterWorkerServiceRequest{}
	if err := json.Unmarshal(requestBytes, &registerWorkerServiceRequest); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}

	if registerWorkerServiceRequest.ServiceId == "" {
		http.Error(w, "service id must not be empty", http.StatusBadRequest)
		return nil
	}

	if registerWorkerServiceRequest.JobName == "" {
		http.Error(w, "job name must not be empty", http.StatusBadRequest)
		return nil
	}

	_, err = s.Queries.CreateService(s.Context, repositories.CreateServiceParams{
		ServiceID: registerWorkerServiceRequest.ServiceId,
		JobName:   sql.NullString{String: registerWorkerServiceRequest.JobName, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error creating service: %w", err)
	}

	w.WriteHeader(200)
	return nil
}

func (s *SchedulerService) UnregisterWorkerService(w http.ResponseWriter, r *http.Request) error {
	serviceId := chi.URLParam(r, "id")

	err := s.Queries.DeleteService(s.Context, serviceId)
	if err != nil {
		return fmt.Errorf("error deleting service: %w", err)
	}

	w.WriteHeader(200)
	return nil
}

func (s *SchedulerService) ScheduleJob(w http.ResponseWriter, r *http.Request) error {
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	scheduleJobRequest := ScheduleJobRequest{}

	if err := json.Unmarshal(requestBytes, &scheduleJobRequest); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}

	err = s.MessageBusService.SendScheduleJobMessage(scheduleJobRequest.JobName, scheduleJobRequest.JobContext)
	if err != nil {
		return err
	}

	return nil
}
