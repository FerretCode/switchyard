package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/caarlos0/env"
	"github.com/ferretcode/switchyard/incident/internal/ingest"
	"github.com/ferretcode/switchyard/incident/internal/prometheus"
	"github.com/ferretcode/switchyard/incident/internal/railway"
	servicemonitor "github.com/ferretcode/switchyard/incident/internal/service_monitor"
	"github.com/ferretcode/switchyard/incident/internal/webhook"
	"github.com/ferretcode/switchyard/incident/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var logger *slog.Logger
var config types.Config
var prometheusCounters types.PrometheusCounters
var incidentStats types.IncidentStats
var deploymentCache types.DeploymentCache

func main() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			logger.Error("error loading environment variables", "err", err)
			return
		}
	}

	if err := env.Parse(&config); err != nil {
		logger.Error("error parsing environment variables", "err", err)
		return
	}

	gqlClient, err := railway.NewClient(&railway.GraphQLClient{
		AuthToken: config.RailwayApiKey,
		BaseURL:   "https://backboard.railway.app/graphql/v2",
	})
	if err != nil {
		logger.Error("error creating graphql client", "err", err)
		return
	}

	done := make(chan bool)
	prometheusCounters = prometheus.Init()
	incidentStats = types.IncidentStats{}
	deploymentCache = types.DeploymentCache{
		Deployments: make(map[string]string),
		Mutex:       sync.Mutex{},
	}

	webhookService := webhook.NewWebhookService(logger, &config)
	ingestService := ingest.NewIngestService(logger, &incidentStats, &config, &prometheusCounters, &webhookService)
	serviceMonitorService := servicemonitor.NewServiceMonitorService(logger, &incidentStats, &config, &prometheusCounters, gqlClient, &deploymentCache, &webhookService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/ingest", func(w http.ResponseWriter, r *http.Request) {
		handleError(ingestService.Ingest(w, r), w, "ingest")
	})

	r.Handle("/metrics", promhttp.Handler())

	go serviceMonitorService.MonitorServices(done)
	go serviceMonitorService.StartCacheCleaningJob(done)

	http.ListenAndServe(":"+config.Port, r)
	done <- true
}

func handleError(err error, w http.ResponseWriter, svc string) {
	if err != nil {
		logger.Error("error processing request", "svc", svc, "err", err)
		http.Error(w, "error processing request", http.StatusInternalServerError)
	}
}
