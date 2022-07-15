package handlers

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/helpers"
	"YP-metrics-and-alerting/internal/models"
	"YP-metrics-and-alerting/internal/render"
	"YP-metrics-and-alerting/internal/repository"
	"YP-metrics-and-alerting/internal/storage"
	"compress/gzip"
	"crypto/hmac"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Repository struct {
	App *config.Application
	DB  repository.DBRepo
}

func NewRepo(appConfig *config.Application, db repository.DBRepo) *Repository {
	return &Repository{
		App: appConfig,
		DB:  db,
	}
}

//UpdateMetricHandler updates metric and returns status
func (repo *Repository) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	metric := &models.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	if metric.MType == models.GaugeType {
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid Value", http.StatusBadRequest)
			return
		}
		metric.Value = &value

	} else if metric.MType == models.CounterType {
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Value", http.StatusBadRequest)
			return
		}

		metric.Delta = &value

	} else {
		http.Error(w, "Unknown metric", http.StatusNotImplemented)
		return
	}

	err := repo.DB.SetMetric(metric)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (repo *Repository) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType != models.GaugeType && metricType != models.CounterType {
		http.Error(w, "Unknown metric", http.StatusBadRequest)
		return
	}

	metric, err := repo.DB.GetMetric(metricName)
	if err != nil {
		http.Error(w, "Metric Not Found", http.StatusBadRequest)
		return
	}

	if metric.MType == models.GaugeType {
		_, err = w.Write([]byte(fmt.Sprintf("%.3f", *metric.Value)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else if metric.MType == models.CounterType {
		_, err = w.Write([]byte(fmt.Sprintf("%d", *metric.Delta)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Metric Type Not Found", http.StatusBadRequest)
	}
}

func (repo *Repository) GetInfoPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	metrics, err := repo.DB.GetMetricsList()
	if err != nil {
		http.Error(w, "Invalid Value", http.StatusBadRequest)
		return
	}

	type htmlPage struct {
		Metrics []*models.Metrics
	}
	data := htmlPage{
		Metrics: metrics,
	}
	err = render.Template(w, r, "metrics.gohtml", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (repo *Repository) UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	var metric models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if metric.MType != models.GaugeType && metric.MType != models.CounterType {
		http.Error(w, "Invalid Metric Type", http.StatusBadRequest)
		return
	}

	if metric.MType == models.GaugeType && metric.Delta != nil {
		http.Error(w, "Use Value instead of Delta", http.StatusBadRequest)
		return
	}

	if metric.MType == models.CounterType && metric.Value != nil {
		http.Error(w, "Use Delta instead of Value", http.StatusBadRequest)
		return
	}

	if key := repo.App.Config.Key; key != "" {
		var str string
		if metric.MType == models.CounterType {
			str = fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta)
		} else if metric.MType == models.GaugeType {
			str = fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)
		}

		hash := helpers.Hash(str, key)
		if !hmac.Equal([]byte(hash), []byte(metric.Hash)) {
			http.Error(w, "Hashes are not equal!", http.StatusBadRequest)
			return
		}
	}

	if repo.App.Config.StoreInterval == 0 {
		repo.SaveMetrics()
	}

	err := repo.DB.SetMetric(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (repo *Repository) UpdateMetricsListJSONHandler(w http.ResponseWriter, r *http.Request) {
	var metrics []*models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if repo.App.Config.StoreInterval == 0 {
		repo.SaveMetrics()
	}

	err := repo.DB.SetMetricsList(metrics)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (repo *Repository) GetMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var metric *models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := repo.DB.GetMetric(metric.ID)
	if err != nil {
		log.Println("GetMetric error:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("%+v\n", metric)

	if key := repo.App.Config.Key; key != "" {
		var str string
		if metric.MType == models.GaugeType {
			str = fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)
		} else if metric.MType == models.CounterType {
			str = fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta)
		}
		log.Println(str)
		metric.Hash = helpers.Hash(str, key)
		log.Println("Hash:", metric.Hash)
	}

	log.Printf("%+v\n", metric)
	if metric.Value != nil {
		log.Println("Value:", *metric.Value)
	}
	if metric.Delta != nil {
		log.Println("Delta:", *metric.Delta)
	}

	err = json.NewEncoder(w).Encode(metric)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (repo *Repository) ServeFileStorage(fileStorage storage.FileStorage) {
	log.Println("ServeFileStorage")

	if repo.App.Config.Restore {
		slice, err := fileStorage.Retrieve()
		if err != nil {
			log.Println(err)
		} else {
			err = repo.DB.SetMetricsListFromFile(slice)
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
		metrics, err := repo.DB.GetMetricsList()
		if err != nil {
			return
		}

		data, err := json.MarshalIndent(metrics, "", "  ")
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

	metricsList, err := repo.DB.GetMetricsList()
	if err != nil {
		return
	}

	data, err := json.MarshalIndent(metricsList, "", "  ")
	if err != nil {
		return
	}

	err = repo.App.FileStorage.Save(data)
	if err != nil {
		return
	}
}

func (repo *Repository) PingDB(w http.ResponseWriter, _ *http.Request) {
	dsn := repo.App.Config.DSN

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	if err = db.Ping(); err != nil {
		log.Printf("Unable to ping database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
