package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/helpers"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

const (
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
)

const (
	Alloc           = "Alloc"
	BuckHashSys     = "BuckHashSys"
	Frees           = "Frees"
	GCCPUFraction   = "GCCPUFraction"
	GCSys           = "GCSys"
	HeapAlloc       = "HeapAlloc"
	HeapIdle        = "HeapIdle"
	HeapInuse       = "HeapInuse"
	HeapObjects     = "HeapObjects"
	HeapReleased    = "HeapReleased"
	HeapSys         = "HeapSys"
	LastGC          = "LastGC"
	Lookups         = "Lookups"
	MCacheInuse     = "MCacheInuse"
	MCacheSys       = "MCacheSys"
	MSpanInuse      = "MSpanInuse"
	MSpanSys        = "MSpanSys"
	Mallocs         = "Mallocs"
	NextGC          = "NextGC"
	NumForcedGC     = "NumForcedGC"
	NumGC           = "NumGC"
	OtherSys        = "OtherSys"
	PauseTotalNs    = "PauseTotalNs"
	StackInuse      = "StackInuse"
	StackSys        = "StackSys"
	Sys             = "Sys"
	PollCount       = "PollCount"
	RandomValue     = "RandomValue"
	TotalAlloc      = "TotalAlloc"
	TotalMemory     = "TotalMemory"
	FreeMemory      = "FreeMemory"
	CPUutilization1 = "CPUutilization1"
)

func InitConfig() *config.AgentConfig {
	cfg := &config.AgentConfig{}

	serverHost := helpers.GetEnv("ADDRESS", "127.0.0.1:8080")
	reportInterval := helpers.StringToSeconds(helpers.GetEnv("REPORT_INTERVAL", "10s"))
	pollInterval := helpers.StringToSeconds(helpers.GetEnv("POLL_INTERVAL", "2s"))
	key := helpers.GetEnv("KEY", "")

	flag.StringVar(&cfg.Address, "a", serverHost, "Send metrics to address:port")
	flag.DurationVar(&cfg.ReportInterval, "r", reportInterval, "Report of interval")
	flag.DurationVar(&cfg.PoolInterval, "p", pollInterval, "Pool of interval")
	flag.StringVar(&cfg.Key, "k", key, "Hashing key")
	flag.Parse()

	protocol := "http"
	cfg.Address = fmt.Sprintf("%s://%s", protocol, cfg.Address)

	return cfg
}

func main() {
	cfg := InitConfig()

	cache := NewAgentCache()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go UpdateMetrics(cfg.PoolInterval, wg, cache)

	wg.Add(1)
	go UpdateExtraMetrics(cfg.PoolInterval, wg, cache)

	wg.Add(1)
	go SendMetrics(cfg, wg, cache)

	//http.ListenAndServe("localhost:8080", nil)

	wg.Wait()
}

func SendMetrics(cfg *config.AgentConfig, wg *sync.WaitGroup, cache *AgentCache) {
	defer wg.Done()

	ticker := time.NewTicker(cfg.ReportInterval)
	for range ticker.C {
		log.Println("Sending metrics...")

		var metricsList = cache.MapToSlice()

		if cfg.Key != "" {
			for _, metric := range metricsList {
				metric.SetHash(cfg.Key)
			}
		}

		url := fmt.Sprintf("%s/updates", cfg.Address)
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

func UpdateExtraMetrics(interval time.Duration, wg *sync.WaitGroup, cache *AgentCache) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	for range ticker.C {
		log.Println("UpdateExtraMetrics...")

		v, err := mem.VirtualMemory()
		if err != nil {
			log.Println(err)
			continue
		}
		cache.Set(TotalMemory, NewGaugeMetric(TotalMemory, float64(v.Total)))
		cache.Set(FreeMemory, NewGaugeMetric(FreeMemory, float64(v.Free)))

		percent, err := cpu.Percent(time.Second, false)
		if err != nil {
			log.Println(err)
			continue
		}
		cache.Set(CPUutilization1, NewGaugeMetric(CPUutilization1, percent[0]))
	}
}

func UpdateMetrics(interval time.Duration, wg *sync.WaitGroup, cache *AgentCache) {
	defer wg.Done()

	var memStats runtime.MemStats
	var pollCount int

	rand.Seed(time.Now().Unix())

	ticker := time.NewTicker(interval)
	for range ticker.C {
		log.Println("Update metrics...")

		runtime.ReadMemStats(&memStats)
		pollCount++

		cache.Set(Alloc, NewGaugeMetric(Alloc, float64(memStats.Alloc)))
		cache.Set(BuckHashSys, NewGaugeMetric(BuckHashSys, float64(memStats.BuckHashSys)))
		cache.Set(Frees, NewGaugeMetric(Frees, float64(memStats.Frees)))
		cache.Set(GCCPUFraction, NewGaugeMetric(GCCPUFraction, memStats.GCCPUFraction))
		cache.Set(GCSys, NewGaugeMetric(GCSys, float64(memStats.GCSys)))
		cache.Set(HeapAlloc, NewGaugeMetric(HeapAlloc, float64(memStats.HeapAlloc)))
		cache.Set(HeapIdle, NewGaugeMetric(HeapIdle, float64(memStats.HeapIdle)))
		cache.Set(HeapInuse, NewGaugeMetric(HeapInuse, float64(memStats.HeapInuse)))
		cache.Set(HeapObjects, NewGaugeMetric(HeapObjects, float64(memStats.HeapObjects)))
		cache.Set(HeapReleased, NewGaugeMetric(HeapReleased, float64(memStats.HeapReleased)))
		cache.Set(HeapSys, NewGaugeMetric(HeapSys, float64(memStats.HeapSys)))
		cache.Set(LastGC, NewGaugeMetric(LastGC, float64(memStats.LastGC)))
		cache.Set(Lookups, NewGaugeMetric(Lookups, float64(memStats.Lookups)))
		cache.Set(MCacheSys, NewGaugeMetric(MCacheSys, float64(memStats.MCacheSys)))
		cache.Set(MCacheInuse, NewGaugeMetric(MCacheInuse, float64(memStats.MCacheInuse)))
		cache.Set(MSpanInuse, NewGaugeMetric(MSpanInuse, float64(memStats.MSpanInuse)))
		cache.Set(MSpanSys, NewGaugeMetric(MSpanSys, float64(memStats.MSpanSys)))
		cache.Set(Mallocs, NewGaugeMetric(Mallocs, float64(memStats.Mallocs)))
		cache.Set(NextGC, NewGaugeMetric(NextGC, float64(memStats.NextGC)))
		cache.Set(NumForcedGC, NewGaugeMetric(NumForcedGC, float64(memStats.NumForcedGC)))
		cache.Set(NumGC, NewGaugeMetric(NumGC, float64(memStats.NumGC)))
		cache.Set(OtherSys, NewGaugeMetric(OtherSys, float64(memStats.OtherSys)))
		cache.Set(PauseTotalNs, NewGaugeMetric(PauseTotalNs, float64(memStats.PauseTotalNs)))
		cache.Set(StackInuse, NewGaugeMetric(StackInuse, float64(memStats.StackInuse)))
		cache.Set(StackSys, NewGaugeMetric(StackSys, float64(memStats.StackSys)))
		cache.Set(TotalAlloc, NewGaugeMetric(TotalAlloc, float64(memStats.TotalAlloc)))
		cache.Set(Sys, NewGaugeMetric(Sys, float64(memStats.Sys)))
		cache.Set(PollCount, NewCounterMetric(PollCount, int64(pollCount)))
		cache.Set(RandomValue, NewGaugeMetric(RandomValue, float64(rand.Intn(10000))))
	}
}
