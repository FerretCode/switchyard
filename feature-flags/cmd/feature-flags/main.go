package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/caarlos0/env"
	"github.com/ferretcode/switchyard/feature-flags/internal/flags"
	"github.com/ferretcode/switchyard/feature-flags/internal/repositories"
	"github.com/ferretcode/switchyard/feature-flags/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var logger *slog.Logger
var config types.Config
var flagStore types.FlagStore

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

	ctx := context.Background()
	queries := repositories.New(conn)

	flagStore := types.FlagStore{
		Mutex: sync.Mutex{},
		Flags: make(map[int]types.Flag),
	}

	flagService := flags.NewFlagsService(logger, &config, &flagStore, queries, ctx)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/flags", func(r chi.Router) {
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			handleError(flagService.Create(w, r), w, "flags/create")
		})

		r.Get("/get/{name}", func(w http.ResponseWriter, r *http.Request) {
			handleError(flagService.Get(w, r), w, "flags/get")
		})

		r.Patch("/update/{name}", func(w http.ResponseWriter, r *http.Request) {
			handleError(flagService.Update(w, r), w, "flags/update")
		})

		r.Delete("/delete/{name}", func(w http.ResponseWriter, r *http.Request) {
			handleError(flagService.Delete(w, r), w, "flags/delete")
		})

		r.Post("/evaluate/{name}", func(w http.ResponseWriter, r *http.Request) {
			handleError(flagService.Evaluate(w, r), w, "flags/evaluate")
		})
	})

	http.ListenAndServe(":"+config.Port, r)
}

func handleError(err error, w http.ResponseWriter, svc string) {
	if err != nil {
		http.Error(w, "there was an error processing your request: "+err.Error(), http.StatusInternalServerError)
		logger.Error("error processing request", "svc", svc, "err", err)
	}
}
