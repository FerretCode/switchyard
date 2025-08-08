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

func (f *FeatureFlagsService) Create(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return nil
	}

	req, err := http.NewRequest("POST", f.Config.FeatureFlagsServiceUrl+"/flags/create", bytes.NewReader(body))
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

func (f *FeatureFlagsService) Get(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "name")

	req, err := http.NewRequest("GET", f.Config.FeatureFlagsServiceUrl+"/flags/get/"+name, nil)
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

func (f *FeatureFlagsService) Update(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "name")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return nil
	}

	req, err := http.NewRequest("PATCH", f.Config.FeatureFlagsServiceUrl+"/flags/update/"+name, bytes.NewReader(body))
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

func (f *FeatureFlagsService) Delete(w http.ResponseWriter, r *http.Request) error {
	name := chi.URLParam(r, "name")

	req, err := http.NewRequest("DELETE", f.Config.FeatureFlagsServiceUrl+"/flags/delete/"+name, nil)
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

func (f *FeatureFlagsService) List(w http.ResponseWriter, r *http.Request) error {
	req, err := http.NewRequest("GET", f.Config.FeatureFlagsServiceUrl+"/flags/list", nil)
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
