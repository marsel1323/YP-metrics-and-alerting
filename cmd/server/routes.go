package main

import (
	"YP-metrics-and-alerting/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(repo *handlers.Repository) *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	//storage := repository.NewMapStorageRepo()
	//repo := handlers.NewRepo(app, storage)

	mux.Get("/", repo.GetAllMetricsHandler)
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", repo.UpdateMetricHandler)
	mux.Get("/value/{metricType}/{metricName}", repo.GetMetricHandler)
	mux.Post("/update", repo.UpdateMetricJSONHandler)
	mux.Post("/update/", repo.UpdateMetricJSONHandler)
	mux.Post("/value", repo.GetMetricJSONHandler)
	mux.Post("/value/", repo.GetMetricJSONHandler)

	return mux
}
