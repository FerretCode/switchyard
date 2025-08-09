package incidentreporting

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/ferretcode/switchyard/dashboard/internal/types"
)

type IncidentReportingService struct {
	Logger *slog.Logger
	Config *types.Config
}

func NewIncidentReportingService(logger *slog.Logger, config *types.Config) IncidentReportingService {
	return IncidentReportingService{
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

func (i *IncidentReportingService) ListIncidentReports(w http.ResponseWriter, r *http.Request) error {
	url := i.Config.IncidentReportingServiceUrl + "/incident/list-incident-reports"
	return PropagateRequest(w, r, "GET", url)
}
