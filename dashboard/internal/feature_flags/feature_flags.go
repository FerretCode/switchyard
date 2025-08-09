package featureflags

import (
	"bytes"
	"io"
	"net/http"

	"log/slog"

	"github.com/ferretcode/switchyard/dashboard/internal/types"
	"github.com/go-chi/chi/v5"
)

type FeatureFlagsService struct {
	Logger *slog.Logger
	Config *types.Config
}

func NewFeatureFlagsService(logger *slog.Logger, config *types.Config) FeatureFlagsService {
	return FeatureFlagsService{
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

func (f *FeatureFlagsService) Create(w http.ResponseWriter, r *http.Request) error {
	url := f.Config.FeatureFlagsServiceUrl + "/flags/create"
	return PropagateRequest(w, r, "POST", url)
}

func (f *FeatureFlagsService) Get(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "name")
	url := f.Config.FeatureFlagsServiceUrl + "/flags/get/" + name
	return PropagateRequest(w, r, "GET", url)
}

func (f *FeatureFlagsService) Update(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "name")
	url := f.Config.FeatureFlagsServiceUrl + "/flags/update/" + name
	return PropagateRequest(w, r, "PATCH", url)
}

func (f *FeatureFlagsService) UpsertRules(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "name")
	url := f.Config.FeatureFlagsServiceUrl + "/flags/upsert-rules/" + name
	return PropagateRequest(w, r, "PATCH", url)
}

func (f *FeatureFlagsService) Delete(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "name")
	url := f.Config.FeatureFlagsServiceUrl + "/flags/delete/" + name
	return PropagateRequest(w, r, "DELETE", url)
}

func (f *FeatureFlagsService) List(w http.ResponseWriter, r *http.Request) error {
	url := f.Config.FeatureFlagsServiceUrl + "/flags/list"
	return PropagateRequest(w, r, "GET", url)
}

func (f *FeatureFlagsService) ToggleFeatureFlag(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "name")
	enabled := r.URL.Query().Get("enabled")
	url := f.Config.FeatureFlagsServiceUrl + "/flags/toggle-feature-flag/" + name + "?enabled=" + enabled
	return PropagateRequest(w, r, "PATCH", url)
}
