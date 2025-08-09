package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/caarlos0/env"
	"github.com/ferretcode/switchyard/dashboard/internal/auth"
	"github.com/ferretcode/switchyard/dashboard/internal/autoscale"
	featureflags "github.com/ferretcode/switchyard/dashboard/internal/feature_flags"
	"github.com/ferretcode/switchyard/dashboard/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

var templates *template.Template
var logger *slog.Logger
var config types.Config

func parseTemplates() error {
	var err error

	files := []string{
		"./views/login.html",
		"./views/scripts.html",
		"./views/header.html",
		"./views/sidebar.html",
		"./views/layout.html",
		"./views/autoscaling.html",
		"./views/feature-flags.html",
		"./views/incident-reporting.html",
		"./views/index.html",
	}

	templates, err = template.ParseFiles(files...)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			logger.Error("error parsing .env", "err", err)
			return
		}
	}

	if err := env.Parse(&config); err != nil {
		logger.Error("error parsing config", "err", err)
		return
	}

	if err := parseTemplates(); err != nil {
		logger.Error("error parsing templates", "err", err)
		return
	}

	authService := auth.NewAuthService(&config)
	autoscaleService := autoscale.NewAutoscaleService(logger, &config)
	featureFlagsService := featureflags.NewFeatureFlagsService(logger, &config)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard/home", http.StatusSeeOther)
	})

	r.Route("/dashboard", func(r chi.Router) {
		// TODO: reenable
		// r.Use(authService.RequireAuth)

		r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
			handleError(templates.ExecuteTemplate(w, "index.html", nil), w, "dashboard/render")
		})
	})

	r.Route("/api", func(r chi.Router) {
		// TODO: consider authentication middleware for dashboard routes

		r.Route("/autoscale", func(r chi.Router) {
			r.Get("/list-services", func(w http.ResponseWriter, r *http.Request) {
				handleError(autoscaleService.ListServices(w, r), w, "autoscale/list")
			})

			r.Patch("/toggle-service-registered", func(w http.ResponseWriter, r *http.Request) {
				handleError(autoscaleService.ToggleServiceRegistered(w, r), w, "autoscale/toggle")
			})
		})

		r.Route("/feature-flags", func(r chi.Router) {
			r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
				handleError(featureFlagsService.Create(w, r), w, "feature-flags/create")
			})

			r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
				handleError(featureFlagsService.List(w, r), w, "feature-flags/list")
			})

			r.Get("/get/{name}", func(w http.ResponseWriter, r *http.Request) {
				handleError(featureFlagsService.Get(w, r), w, "feature-flags/get")
			})

			r.Patch("/update/{name}", func(w http.ResponseWriter, r *http.Request) {
				handleError(featureFlagsService.Update(w, r), w, "feature-flags/update")
			})

			r.Patch("/upsert-rules/{name}", func(w http.ResponseWriter, r *http.Request) {
				handleError(featureFlagsService.UpsertRules(w, r), w, "feature-flags/upsert-rules")
			})

			r.Patch("/toggle-feature-flag/{name}", func(w http.ResponseWriter, r *http.Request) {
				handleError(featureFlagsService.ToggleFeatureFlag(w, r), w, "feature-flags/enable-feature-flag")
			})

			r.Delete("/delete/{name}", func(w http.ResponseWriter, r *http.Request) {
				handleError(featureFlagsService.Delete(w, r), w, "feature-flags/delete")
			})
		})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			handleError(authService.RenderLogin(w, r, templates), w, "login/render")
		})

		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			handleError(authService.Login(w, r), w, "login")
		})

		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			handleError(authService.Logout(w, r), w, "logout")
		})
	})

	http.ListenAndServe(":"+config.Port, r)
}

func handleError(err error, w http.ResponseWriter, svc string) {
	if err != nil {
		http.Error(w, "there was an error processing your request", http.StatusInternalServerError)
		logger.Error("error processing request", "svc", svc, "err", err)
	}
}
