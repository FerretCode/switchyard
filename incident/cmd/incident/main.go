package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/caarlos0/env"
	"github.com/ferretcode/switchyard/incident/internal/api"
	"github.com/ferretcode/switchyard/incident/internal/ingest"
	messagebus "github.com/ferretcode/switchyard/incident/internal/message_bus"
	"github.com/ferretcode/switchyard/incident/internal/prometheus"
	"github.com/ferretcode/switchyard/incident/internal/railway"
	"github.com/ferretcode/switchyard/incident/internal/repositories"
	servicemonitor "github.com/ferretcode/switchyard/incident/internal/service_monitor"
	"github.com/ferretcode/switchyard/incident/internal/webhook"
	"github.com/ferretcode/switchyard/incident/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"

	_ "github.com/lib/pq"
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

	conn, err := sqlx.Open("postgres", config.DatabaseUrl)
	if err != nil {
		logger.Error("error opening database connection", "err", err)
		return
	}
	defer conn.Close()

	messageBusConn, err := amqp.Dial(config.MessageBusUrl)
	if err != nil {
		logger.Error("error connecting to the message bus", "err", err)
	}
	defer messageBusConn.Close()

	ctx := context.Background()
	queries := repositories.New(conn)

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

	messageBusService := messagebus.NewMessageBusService(logger, messageBusConn, &config, ctx)
	webhookService := webhook.NewWebhookService(logger, &config, queries, ctx, &messageBusService)
	ingestService := ingest.NewIngestService(logger, &incidentStats, &config, &prometheusCounters, &webhookService)
	serviceMonitorService := servicemonitor.NewServiceMonitorService(logger, &incidentStats, &config, &prometheusCounters, gqlClient, &deploymentCache, &webhookService)
	apiService := api.NewApiService(logger, &config, queries, ctx)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/incident", func(r chi.Router) {
		r.Post("/ingest", func(w http.ResponseWriter, r *http.Request) {
			handleError(ingestService.Ingest(w, r), w, "ingest")
		})

		r.Get("/list-incident-reports", func(w http.ResponseWriter, r *http.Request) {
			handleError(apiService.ListIncidentReports(w, r), w, "list-incident-reports")
		})
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
