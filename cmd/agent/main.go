package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/helpers"
	"YP-metrics-and-alerting/internal/models"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
)

const (
	Alloc         = "Alloc"
	BuckHashSys   = "BuckHashSys"
	Frees         = "Frees"
	GCCPUFraction = "GCCPUFraction"
	GCSys         = "GCSys"
	HeapAlloc     = "HeapAlloc"
	HeapIdle      = "HeapIdle"
	HeapInuse     = "HeapInuse"
	HeapObjects   = "HeapObjects"
	HeapReleased  = "HeapReleased"
	HeapSys       = "HeapSys"
	LastGC        = "LastGC"
	Lookups       = "Lookups"
	MCacheInuse   = "MCacheInuse"
	MCacheSys     = "MCacheSys"
	MSpanInuse    = "MSpanInuse"
	MSpanSys      = "MSpanSys"
	Mallocs       = "Mallocs"
	NextGC        = "NextGC"
	NumForcedGC   = "NumForcedGC"
	NumGC         = "NumGC"
	OtherSys      = "OtherSys"
	PauseTotalNs  = "PauseTotalNs"
	StackInuse    = "StackInuse"
	StackSys      = "StackSys"
	Sys           = "Sys"
	PollCount     = "PollCount"
	RandomValue   = "RandomValue"
	TotalAlloc    = "TotalAlloc"
)

func main() {
	cfg := &config.AgentConfig{}

	serverHost := helpers.GetEnv("ADDRESS", "127.0.0.1:8080")
	reportInterval := helpers.StringToSeconds(helpers.GetEnv("REPORT_INTERVAL", "10s"))
	pollInterval := helpers.StringToSeconds(helpers.GetEnv("POLL_INTERVAL", "2s"))
	key := helpers.GetEnv("KEY", "secretkey")

	flag.StringVar(&cfg.Address, "a", serverHost, "Send metrics to address:port")
	flag.DurationVar(&cfg.ReportInterval, "r", reportInterval, "Report of interval")
	flag.DurationVar(&cfg.PoolInterval, "p", pollInterval, "Pool of interval")
	flag.StringVar(&cfg.Key, "k", key, "Hashing key")
	flag.Parse()

	protocol := "http"
	cfg.Address = fmt.Sprintf("%s://%s", protocol, cfg.Address)

	log.Println(cfg)

	metricsMap := NewMetricsMap()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go metricsMap.UpdateMetrics(cfg.PoolInterval, wg)

	wg.Add(1)
	go metricsMap.SendMetrics(cfg, wg)

	wg.Wait()
}

type MetricsMap map[string]interface{}

func NewMetricsMap() MetricsMap {
	return make(map[string]interface{})
}

func (metricsMap MetricsMap) SendMetrics(cfg *config.AgentConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	serverHost := cfg.Address
	interval := cfg.ReportInterval
	key := cfg.Key

	ticker := time.NewTicker(interval)
	for range ticker.C {
		log.Println("Sending metrics...")

		var metricsList []*models.Metrics

		for metricName, value := range metricsMap {
			metricType := GaugeMetricType
			if metricName == PollCount {
				metricType = CounterMetricType
			}

			metric := &models.Metrics{
				ID:    metricName,
				MType: metricType,
			}

			if metricType == CounterMetricType {
				metricValue := value.(int64)
				metric.Delta = &metricValue
			} else if metricType == GaugeMetricType {
				metricValue := value.(float64)
				metric.Value = &metricValue
			}

			if key != "" {
				if metricType == CounterMetricType {
					src := fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta)
					log.Println(src)
					metric.Hash = helpers.Hash(src, key)
				}

				if metricType == GaugeMetricType {
					src := fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)
					log.Println(src)
					metric.Hash = helpers.Hash(src, key)
				}
			}

			log.Printf("%+v\n", metric)
			metricsList = append(metricsList, metric)
		}

		url := fmt.Sprintf("%s/updates", serverHost)
		body, err := json.Marshal(metricsList)
		if err != nil {
			log.Println(err)
			return
		}
		request, err := http.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			log.Println("Unable to send metrics to server:", err)
			continue
		}

		err = request.Body.Close()
		if err != nil {
			log.Println(err)
			break
		}
	}

}

func (metricsMap MetricsMap) UpdateMetrics(interval time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	var memStats runtime.MemStats
	var pollCount int

	rand.Seed(time.Now().Unix())

	ticker := time.NewTicker(interval)
	for range ticker.C {
		log.Println("Updating metrics...")

		pollCount++

		runtime.ReadMemStats(&memStats)

		metricsMap[Alloc] = float64(memStats.Alloc)
		metricsMap[BuckHashSys] = float64(memStats.BuckHashSys)
		metricsMap[Frees] = float64(memStats.Frees)
		metricsMap[GCCPUFraction] = memStats.GCCPUFraction
		metricsMap[GCSys] = float64(memStats.GCSys)
		metricsMap[HeapAlloc] = float64(memStats.HeapAlloc)
		metricsMap[HeapIdle] = float64(memStats.HeapIdle)
		metricsMap[HeapInuse] = float64(memStats.HeapInuse)
		metricsMap[HeapObjects] = float64(memStats.HeapObjects)
		metricsMap[HeapReleased] = float64(memStats.HeapReleased)
		metricsMap[HeapSys] = float64(memStats.HeapSys)
		metricsMap[LastGC] = float64(memStats.LastGC)
		metricsMap[Lookups] = float64(memStats.Lookups)
		metricsMap[MCacheInuse] = float64(memStats.MCacheInuse)
		metricsMap[MCacheSys] = float64(memStats.MCacheSys)
		metricsMap[MSpanInuse] = float64(memStats.MSpanInuse)
		metricsMap[MSpanSys] = float64(memStats.MSpanSys)
		metricsMap[Mallocs] = float64(memStats.Mallocs)
		metricsMap[NextGC] = float64(memStats.NextGC)
		metricsMap[NumForcedGC] = float64(memStats.NumForcedGC)
		metricsMap[NumGC] = float64(memStats.NumGC)
		metricsMap[OtherSys] = float64(memStats.OtherSys)
		metricsMap[PauseTotalNs] = float64(memStats.PauseTotalNs)
		metricsMap[StackInuse] = float64(memStats.StackInuse)
		metricsMap[StackSys] = float64(memStats.StackSys)
		metricsMap[TotalAlloc] = float64(memStats.TotalAlloc)
		metricsMap[Sys] = float64(memStats.Sys)
		metricsMap[PollCount] = int64(pollCount)
		metricsMap[RandomValue] = float64(rand.Intn(10000))
	}
}
