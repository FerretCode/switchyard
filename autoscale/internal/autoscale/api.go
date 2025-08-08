package autoscale

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ferretcode/switchyard/autoscale/internal/repositories"
	"github.com/go-chi/chi/v5"
)

func (a *AutoscaleService) RegisterService(w http.ResponseWriter, r *http.Request) error {
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading reques body: %w", err)
	}

	registerServiceRequest := RegisterServiceRequest{}

	if err := json.Unmarshal(requestBytes, &registerServiceRequest); err != nil {
		return fmt.Errorf("error parsing request body: %w", err)
	}

	if registerServiceRequest.ServiceId == "" {
		http.Error(w, "service id must not be empty", http.StatusBadRequest)
		return nil
	}

	_, err = a.Queries.CreateService(a.Context, repositories.CreateServiceParams{
		ServiceID: registerServiceRequest.ServiceId,
		JobName: sql.NullString{
			String: registerServiceRequest.JobName,
			Valid:  true,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating service: %w", err)
	}

	w.WriteHeader(200)
	return nil
}

func (a *AutoscaleService) UnregisterService(w http.ResponseWriter, r *http.Request) error {
	serviceId := chi.URLParam(r, "id")

	err := a.Queries.DeleteService(a.Context, serviceId)
	if err != nil {
		return fmt.Errorf("error deleting service: %w", err)
	}

	w.WriteHeader(200)
	return nil
}
