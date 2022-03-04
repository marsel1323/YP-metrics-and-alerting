package handlers

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/models"
	"YP-metrics-and-alerting/internal/repository"
	"YP-metrics-and-alerting/internal/storage"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	log.Println("UpdateMetricHandler")
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
	log.Println("GetMetricHandler")
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	log.Println(metricType, metricName)

	if metricType == GaugeType {
		value, err := repo.DB.GetGaugeMetricValue(metricName)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}
		log.Println(value)
		_, err = w.Write([]byte(fmt.Sprintf("%.3f", value)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	} else if metricType == CounterType {
		value, err := repo.DB.GetCounterMetricValue(metricName)
		if err != nil {
			http.Error(w, "Metric Not Found", http.StatusNotFound)
			return
		}
		log.Println(value)
		_, err = w.Write([]byte(fmt.Sprintf("%d", value)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "Metric Type Not Found", http.StatusNotFound)
}

func (repo *Repository) GetInfoPageHandler(w http.ResponseWriter, _ *http.Request) {
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

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (repo *Repository) UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateMetricJSONHandler")

	var m models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metricType := m.MType
	metricName := m.ID

	log.Println(m.MType, m.ID)
	if m.Value != nil {
		log.Println("m.Value", *m.Value)
	}

	if m.Delta != nil {
		log.Println("m.Delta", *m.Delta)
	}

	if repo.App.Config.StoreInterval == 0 {
		log.Println("StoreInterval == 0")
		repo.SaveMetrics()
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
	log.Println("GetMetricJSONHandler")

	w.Header().Set("Content-Type", "application/json")

	var m models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricType := m.MType
	metricName := m.ID

	log.Println(metricType, metricName)

	if metricType == GaugeType {
		handleGaugeMetric(w, &m, repo)
	} else if metricType == CounterType {
		handleCounterMetric(w, &m, repo)
	} else {
		log.Println("Metric Type Not Found")
		http.Error(w, "Metric Type Not Found", http.StatusNotFound)
	}
}

func handleCounterMetric(w http.ResponseWriter, m *models.Metrics, repo *Repository) {
	value, err := repo.DB.GetCounterMetricValue(m.ID)
	if err != nil {
		log.Println(err)
	}

	m.Delta = &value

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func handleGaugeMetric(w http.ResponseWriter, m *models.Metrics, repo *Repository) {
	value, err := repo.DB.GetGaugeMetricValue(m.ID)
	if err != nil {
		log.Println(err)
	}

	m.Value = &value

	err = json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (repo *Repository) ServeFileStorage(fileStorage storage.Storage) {
	log.Println("ServeFileStorage")

	if repo.App.Config.Restore {
		mapStorage, err := fileStorage.Retrieve()
		if err != nil {
			log.Println(err)
		} else {
			err = repo.DB.BunchSetMetrics(mapStorage)
			if err != nil {
				log.Println(err)
			}
		}
	}

	if repo.App.Config.StoreInterval == 0 {
		log.Println("STORE_INTERVAL == 0")
		return
	}

	storeTickerInterval := time.NewTicker(repo.App.Config.StoreInterval)
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

func (repo *Repository) SaveMetrics() {
	log.Println("SaveMetrics")

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
	if err != nil {
		return
	}
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type ContentType []string

var contentTypes ContentType

func init() {
	contentTypes = append(contentTypes, "application/javascript")
	contentTypes = append(contentTypes, "application/json")
	contentTypes = append(contentTypes, "text/css")
	contentTypes = append(contentTypes, "text/html")
	contentTypes = append(contentTypes, "text/plain")
	contentTypes = append(contentTypes, "text/xml")
}

func (c *ContentType) Contains(value string) bool {
	for _, ct := range contentTypes {
		if strings.Contains(value, ct) {
			return true
		}
	}
	return false
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		log.Println(contentType)
		log.Println(contentTypes.Contains(contentType))
		if !contentTypes.Contains(contentType) {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, err = io.WriteString(w, err.Error())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
