package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/handlers"
	"YP-metrics-and-alerting/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(app *config.Application) *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	storage := repository.NewMapStorageRepo()
	repo := handlers.NewRepo(app, storage)

	mux.Get("/", repo.GetAllMetricsHandler)
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", repo.UpdateMetricHandler)
	mux.Get("/value/{metricType}/{metricName}", repo.GetMetricHandler)

	return mux
}