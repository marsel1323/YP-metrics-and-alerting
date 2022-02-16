package handlers

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/repository"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
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
	log.Println(metricType, metricName, metricValue)

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
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType == GaugeType {
		value, err := repo.DB.GetGaugeMetricValue(metricName)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}

		w.Write([]byte(fmt.Sprintf("%f", value)))
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

	data := make(map[string]interface{})
	data["GaugeMetrics"] = gaugeMetrics
	data["CounterMetrics"] = counterMetrics

	t, err := template.New("metricsHTML").Parse(`
		<html>
			<head>
				<title>Metrics</title>
				<meta http-equiv="refresh" content="10" />
		  	</head>
			<body>
				<h2>Gauge Metrics</h1>
				{{ $gaugeMetrics := index . "GaugeMetrics" }}
				<ol>
					{{ range $key, $value := $gaugeMetrics }}
						<li>
							<b>{{ $key }}</b>: {{ $value }}
						</li>
					{{ end }}
				</ol>
				
				<h2>Counter Metrics</h1>
				{{ $counterMetrics := index . "CounterMetrics" }}
				<ol>
					{{ range $key, $value := $counterMetrics }}
						<li>
							<b>{{ $key }}</b>: {{ $value }}
						</li>
					{{ end }}
				</ol>
			</body>
		</html>
		`)
	if err != nil {
		log.Println(err)
	}

	t.Execute(w, data)
}
