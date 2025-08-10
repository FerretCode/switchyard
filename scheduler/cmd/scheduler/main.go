package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/caarlos0/env/v10"
	messagebus "github.com/ferretcode/switchyard/scheduler/internal/message_bus"
	"github.com/ferretcode/switchyard/scheduler/internal/repositories"
	"github.com/ferretcode/switchyard/scheduler/internal/scheduler"
	"github.com/ferretcode/switchyard/scheduler/internal/watchdog"
	"github.com/ferretcode/switchyard/scheduler/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var logger *slog.Logger
var config types.Config

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

	options, err := redis.ParseURL(config.CacheUrl)
	if err != nil {
		logger.Error("error parsing redis url", "err", err)
		return
	}

	redisConn := redis.NewClient(options)

	messageBusConn, err := amqp.Dial(config.MessageBusUrl)
	if err != nil {
		logger.Error("error connecting to the message bus", "err", err)
	}
	defer messageBusConn.Close()

	queries := repositories.New(conn)

	messageBusService := messagebus.NewMessageBusService(logger, messageBusConn, &config, redisConn, ctx, queries)
	watchdogService := watchdog.NewWatchdogService(logger, &config, redisConn, &messageBusService, ctx, queries)
	schedulerService := scheduler.NewSchedulerService(logger, queries, ctx, &messageBusService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/scheduler", func(r chi.Router) {
		r.Post("/register-worker-service", func(w http.ResponseWriter, r *http.Request) {
			handleError(schedulerService.RegisterWorkerService(w, r), w, "scheduler/register-worker-service")
		})

		r.Delete("/unregister-worker-service/{id}", func(w http.ResponseWriter, r *http.Request) {
			handleError(schedulerService.UnregisterWorkerService(w, r), w, "scheduler/register-worker-service")
		})

		r.Post("/schedule-job", func(w http.ResponseWriter, r *http.Request) {
			handleError(schedulerService.ScheduleJob(w, r), w, "scheduler/schedule-job")
		})

		r.Get("/get-job-statistics/{name}", func(w http.ResponseWriter, r *http.Request) {
			handleError(schedulerService.GetJobStatistics(w, r), w, "scheduler/get-job-statistics")
		})
	})

	go messageBusService.SubscribeToJobFinishedMessages()
	go watchdogService.WatchStuckJobs()

	http.ListenAndServe(":"+config.Port, r)
}

func handleError(err error, w http.ResponseWriter, svc string) {
	if err != nil {
		http.Error(w, "there was an error processing your request: "+err.Error(), http.StatusInternalServerError)
		logger.Error("error processing request", "svc", svc, "err", err)
	}
}
