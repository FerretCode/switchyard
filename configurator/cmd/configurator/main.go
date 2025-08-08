package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/caarlos0/env"
	"github.com/ferretcode/switchyard/configurator/internal/configurator"
	"github.com/ferretcode/switchyard/configurator/internal/railway"
	"github.com/ferretcode/switchyard/configurator/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
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

	gqlClient, err := railway.NewClient(&railway.GraphQLClient{
		AuthToken: config.RailwayApiKey,
		BaseURL:   "https://backboard.railway.app/graphql/v2",
	})
	if err != nil {
		logger.Error("error creating graphql client", "err", err)
		return
	}

	configuratorService := configurator.NewConfiguratorService(logger, &config, gqlClient, ctx)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/configurator", func(r chi.Router) {
		r.Post("/configure/{service}", func(w http.ResponseWriter, r *http.Request) {
			handleError(configuratorService.UpdateConfig(w, r), w, "configurator/configure")
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
