package scheduler

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/ferretcode/switchyard/dashboard/internal/types"
	"github.com/go-chi/chi/v5"
)

type SchedulerService struct {
	Logger *slog.Logger
	Config *types.Config
}

func NewSchedulerService(logger *slog.Logger, config *types.Config) SchedulerService {
	return SchedulerService{
		Logger: logger,
		Config: config,
	}
}

func PropagateRequest(w http.ResponseWriter, r *http.Request, method, url string) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return nil
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	w.WriteHeader(res.StatusCode)
	w.Write(responseBytes)
	return nil
}

func (s *SchedulerService) RegisterWorkerService(w http.ResponseWriter, r *http.Request) error {
	url := s.Config.AutoscaleServiceUrl + "/autoscale/upsert-service"
	return PropagateRequest(w, r, "POST", url)
}

func (s *SchedulerService) ScheduleJob(w http.ResponseWriter, r *http.Request) error {
	url := s.Config.SchedulerServiceUrl + "/scheduler/schedule-job"
	return PropagateRequest(w, r, "POST", url)
}

func (s *SchedulerService) GetJobStatistics(w http.ResponseWriter, r *http.Request) error {
	url := s.Config.SchedulerServiceUrl + "/scheduler/get-job-statistics/" + chi.URLParam(r, "name")
	return PropagateRequest(w, r, "GET", url)
}

func (s *SchedulerService) UnregisterWorkerService(w http.ResponseWriter, r *http.Request) error {
	url := s.Config.SchedulerServiceUrl + "/scheduler/unregister-worker-service"
	return PropagateRequest(w, r, "DELETE", url)
}
