package main

import (
	"YP-metrics-and-alerting/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "net/http/pprof"
)

func Routes(repo *handlers.Repository) *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)

	mux.Get("/", repo.GetInfoPageHandler)
	mux.Post("/update/{metricType}/{metricName}/{metricValue}", repo.UpdateMetricHandler)
	mux.Get("/value/{metricType}/{metricName}", repo.GetMetricHandler)
	mux.Post("/update*", repo.UpdateMetricJSONHandler)
	//mux.Post("/update/", repo.UpdateMetricJSONHandler)
	mux.Post("/updates*", repo.UpdateMetricsListJSONHandler)
	mux.Post("/value*", repo.GetMetricJSONHandler)
	//mux.Post("/value/", repo.GetMetricJSONHandler)
	mux.Get("/ping*", repo.PingDB)
	mux.Mount("/debug", middleware.Profiler())

	return mux
}
