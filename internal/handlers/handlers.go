package handlers

import (
	"YP-metrics-and-alerting/internal/models"
	"net/http"
	"strconv"
	"strings"
)

type Storage struct {
	Gauge map[string]*models.Gauge
}

type GaugeMetric struct {
	Value float64
}

var gaugeStorage = make(map[string]*GaugeMetric)

type CounterMetric struct {
	Value int64
}

var counterStorage = make(map[string]*CounterMetric)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.RequestURI()

	splitURI := strings.Split(q, "/")
	if len(splitURI) < 5 {
		http.Error(w, "URL Not Found", http.StatusNotFound)
		return
	}

	metricType, metricName, metricValue := splitURI[2], splitURI[3], splitURI[4]

	if metricType == "gauge" {
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid Value", http.StatusBadRequest)
			return
		}
		gaugeStorage[metricName] = &GaugeMetric{
			Value: value,
		}
	} else if metricType == "counter" {
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Value", http.StatusBadRequest)
			return
		}
		counterStorage[metricName] = &CounterMetric{
			Value: value,
		}
	} else {
		http.Error(w, "Unknown metric", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func StatusHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))
}
