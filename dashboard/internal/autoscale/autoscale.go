package autoscale

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/ferretcode/switchyard/dashboard/internal/types"
)

type AutoscaleService struct {
	Logger *slog.Logger
	Config *types.Config
}

func NewAutoscaleService(logger *slog.Logger, config *types.Config) AutoscaleService {
	return AutoscaleService{
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

func (a *AutoscaleService) ConfigureService(w http.ResponseWriter, r *http.Request) error {
	url := a.Config.AutoscaleServiceUrl + "/autoscale/upsert-service"
	return PropagateRequest(w, r, "POST", url)
}

func (a *AutoscaleService) ToggleServiceRegistered(w http.ResponseWriter, r *http.Request) error {
	serviceId := r.URL.Query().Get("service")
	if serviceId == "" {
		http.Error(w, "Service id is required", http.StatusBadRequest)
		return nil
	}

	enabled := r.URL.Query().Get("enabled")
	if serviceId == "" {
		http.Error(w, "Enabled flag is required", http.StatusBadRequest)
		return nil
	}

	var err error
	var req *http.Request

	if enabled == "true" {
		upsertServiceRequest := UpsertServiceRequest{
			ServiceId: &serviceId,
			Enabled:   newBoolPtr(true),
		}
		requestBytes, err := json.Marshal(upsertServiceRequest)
		if err != nil {
			http.Error(w, "Error creating request body", http.StatusInternalServerError)
			return nil
		}

		req, err = http.NewRequest(
			"POST", a.Config.AutoscaleServiceUrl+"/autoscale/upsert-service", bytes.NewReader(requestBytes),
		)
		if err != nil {
			return err
		}
	} else {
		req, err = http.NewRequest(
			"DELETE", a.Config.AutoscaleServiceUrl+"/autoscale/unregister-service/"+serviceId, nil,
		)
		if err != nil {
			return err
		}
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

func (a *AutoscaleService) ListServices(w http.ResponseWriter, r *http.Request) error {
	req, err := http.NewRequest(
		"GET", a.Config.AutoscaleServiceUrl+"/autoscale/list-services", nil,
	)
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

	listServicesResponse := ListServicesResponse{}

	if err := json.Unmarshal(responseBytes, &listServicesResponse); err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(responseBytes) // we'll just send back what we got from the server originally

	return nil
}

func newBoolPtr(value bool) *bool {
	b := value
	return &b
}
