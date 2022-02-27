package handlers

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/models"
	"YP-metrics-and-alerting/internal/repository"
	"YP-metrics-and-alerting/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type Repository struct {
	App *config.Application
	DB  repository.StorageRepo
}

func NewRepo(appConfig *config.Application, db repository.StorageRepo) *Repository {
	return &Repository{
		App: appConfig,
		DB:  db,
	}
}

func (repo *Repository) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricType == GaugeType {
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid Value", http.StatusBadRequest)
			return
		}

		err = repo.DB.SetGaugeMetricValue(metricName, value)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	if metricType == CounterType {
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Value", http.StatusBadRequest)
			return
		}

		err = repo.DB.SetCounterMetricValue(metricName, value)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	if repo.App.StoreInterval == 0 {
		log.Println("StoreInterval == 0")
		repo.Jsonchik()
	}

	http.Error(w, "Unknown metric", http.StatusNotImplemented)
}

func (repo *Repository) Jsonchik() {
	log.Println("Jsonchik")
	gaugeMetrics, err := repo.DB.GetAllGaugeMetricValues()
	if err != nil {
		return
	}

	counterMetrics, err := repo.DB.GetAllCounterMetricValues()
	if err != nil {
		return
	}

	data, err := json.MarshalIndent(repository.MapStorageRepo{
		Gauge:   gaugeMetrics,
		Counter: counterMetrics,
	}, "", "  ")
	if err != nil {
		return
	}

	err = repo.App.FileStorage.Save(data)
}

func (repo *Repository) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType == GaugeType {
		value, err := repo.DB.GetGaugeMetricValue(metricName)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}

		w.Write([]byte(fmt.Sprintf("%.3f", value)))
		return
	} else if metricType == CounterType {
		value, err := repo.DB.GetCounterMetricValue(metricName)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}

		w.Write([]byte(fmt.Sprintf("%d", value)))
		return
	} else {
		http.Error(w, "Metric Type Not Found", http.StatusNotFound)
	}
}

func (repo *Repository) GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	gaugeMetrics, err := repo.DB.GetAllGaugeMetricValues()
	if err != nil {
		http.Error(w, "Invalid Value", http.StatusBadRequest)
		return
	}

	counterMetrics, err := repo.DB.GetAllCounterMetricValues()
	if err != nil {
		http.Error(w, "Invalid Value", http.StatusBadRequest)
		return
	}

	type htmlPage struct {
		GaugeMetrics   map[string]float64
		CounterMetrics map[string]int64
	}

	data := htmlPage{
		GaugeMetrics:   gaugeMetrics,
		CounterMetrics: counterMetrics,
	}

	t, err := template.ParseFiles("../../internal/templates/metrics.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func (repo *Repository) UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	var m models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricType := m.MType
	metricName := m.ID

	log.Println("StoreInterval", repo.App.StoreInterval)
	if repo.App.StoreInterval == 0 {
		log.Println("StoreInterval == 0")
		repo.Jsonchik()
	}

	if metricType == GaugeType {
		metricValue := *m.Value

		err := repo.DB.SetGaugeMetricValue(metricName, metricValue)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	if metricType == CounterType {
		metricValue := *m.Delta

		err := repo.DB.SetCounterMetricValue(metricName, metricValue)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "Unknown metric", http.StatusNotImplemented)
}

func (repo *Repository) GetMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricType := m.MType
	metricName := m.ID

	if metricType == GaugeType {
		value, err := repo.DB.GetGaugeMetricValue(metricName)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}

		m.Value = &value

		err = json.NewEncoder(w).Encode(m)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}

		return
	} else if metricType == CounterType {
		value, err := repo.DB.GetCounterMetricValue(metricName)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}

		m.Delta = &value

		err = json.NewEncoder(w).Encode(m)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}

		return
	} else {
		http.Error(w, "Metric Type Not Found", http.StatusNotFound)
	}
}

func (repo *Repository) ServeFileStorage(fileStorage storage.Storage) {
	if repo.App.Restore {
		err := fileStorage.Retrieve()
		if err != nil {
			log.Println(err)
		}
	}

	if repo.App.StoreInterval == 0 {
		log.Println("STORE_INTERVAL == 0")
		return
	}

	storeTickerInterval := time.NewTicker(repo.App.StoreInterval)
	for range storeTickerInterval.C {
		gaugeMetrics, err := repo.DB.GetAllGaugeMetricValues()
		if err != nil {
			return
		}

		counterMetrics, err := repo.DB.GetAllCounterMetricValues()
		if err != nil {
			return
		}

		data, err := json.MarshalIndent(repository.MapStorageRepo{
			Gauge:   gaugeMetrics,
			Counter: counterMetrics,
		}, "", "  ")
		if err != nil {
			return
		}

		err = fileStorage.Save(data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
