package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/ferretcode/switchyard/autoscale/internal/autoscale"
	"github.com/ferretcode/switchyard/autoscale/internal/railway"
	"github.com/ferretcode/switchyard/autoscale/internal/repositories"
	"github.com/ferretcode/switchyard/autoscale/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var logger *slog.Logger
var config types.Config

const (
	monitoringInterval = 10 * time.Second
	metricHistorySize  = 12
	upscaleCooldown    = 2 * time.Minute
	downscaleCooldown  = 1 * time.Minute
	minReplicaCount    = 1
	maxReplicaCount    = 10

	spikeThreshold = 0.7
	spikeWindow    = 3
)

var (
	lastUpscaleTime     time.Time
	lastDownscaleTime   time.Time
	consecutiveHighLoad int
	consecutiveLowLoad  int
)

func main() {
	ctx := context.Background()

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			logger.Error("error loading .env", "err", err)
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

	gqlClient, err := railway.NewClient(&railway.GraphQLClient{
		AuthToken: config.RailwayApiKey,
		BaseURL:   "https://backboard.railway.app/graphql/v2",
	})
	if err != nil {
		logger.Error("error creating graphql client", "err", err)
		return
	}

	queries := repositories.New(conn)
	serviceHistoryCache := types.ServiceHistoryCache{
		ServiceHistories: make(map[string]types.MetricHistory),
	}

	gqlQueries := railway.NewQueryService(gqlClient, ctx, config, logger)
	autoscalingService := autoscale.NewAutoscaleService(logger, &config, &gqlQueries, queries, ctx, &serviceHistoryCache)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/autoscale", func(r chi.Router) {
		r.Post("/register-service", func(w http.ResponseWriter, r *http.Request) {
			handleError(autoscalingService.RegisterService(w, r), w, "autoscale/register")
		})

		r.Delete("/unregister-service/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleError(autoscalingService.UnregisterService(w, r), w, "autoscale/unregister")
		})
	})

	go autoscalingService.StartAutoscaling()
	http.ListenAndServe(":"+config.Port, r)
}

func handleError(err error, w http.ResponseWriter, svc string) {
	if err != nil {
		http.Error(w, "there was an error processing your request: "+err.Error(), http.StatusInternalServerError)
		logger.Error("error processing request", "svc", svc, "err", err)
	}
}
