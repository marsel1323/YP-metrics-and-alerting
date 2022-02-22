package handlers

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/repository"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
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

	http.Error(w, "Unknown metric", http.StatusNotImplemented)
}

func (repo *Repository) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	var m Metrics

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	m.MType = metricType
	m.ID = metricName

	if metricType == GaugeType {
		value, err := repo.DB.GetGaugeMetricValue(metricName)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}

		//w.Write([]byte(fmt.Sprintf("%.3f", value)))

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

		//w.Write([]byte(fmt.Sprintf("%d", value)))
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

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (repo *Repository) UpdateMetricJsonHandler(w http.ResponseWriter, r *http.Request) {
	var m Metrics

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricType := m.MType
	metricName := m.ID

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

func (repo *Repository) GetMetricJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m Metrics

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
