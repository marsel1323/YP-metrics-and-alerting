package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Gauge float64
type Counter int64

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
)

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	host := "127.0.0.1"
	port := 8080

	metricsMap := make(map[string]interface{})

	go func(interval time.Duration) {
		var memStats runtime.MemStats
		var pollCount int

		rand.Seed(time.Now().Unix())

		for {
			<-time.After(interval)

			log.Println("Updating metrics...")

			pollCount++

			runtime.ReadMemStats(&memStats)

			metricsMap[Alloc] = Gauge(memStats.Alloc)
			metricsMap[BuckHashSys] = Gauge(memStats.BuckHashSys)
			metricsMap[Frees] = Gauge(memStats.Frees)
			metricsMap[GCCPUFraction] = memStats.GCCPUFraction
			metricsMap[GCSys] = Gauge(memStats.GCSys)
			metricsMap[HeapAlloc] = Gauge(memStats.HeapAlloc)
			metricsMap[HeapIdle] = Gauge(memStats.HeapIdle)
			metricsMap[HeapInuse] = Gauge(memStats.HeapInuse)
			metricsMap[HeapObjects] = Gauge(memStats.HeapObjects)
			metricsMap[HeapReleased] = Gauge(memStats.HeapReleased)
			metricsMap[HeapSys] = Gauge(memStats.HeapSys)
			metricsMap[LastGC] = Gauge(memStats.LastGC)
			metricsMap[Lookups] = Gauge(memStats.Lookups)
			metricsMap[MCacheInuse] = Gauge(memStats.MCacheInuse)
			metricsMap[MCacheSys] = Gauge(memStats.MCacheSys)
			metricsMap[MSpanInuse] = Gauge(memStats.MSpanInuse)
			metricsMap[MSpanSys] = Gauge(memStats.MSpanSys)
			metricsMap[Mallocs] = Gauge(memStats.Mallocs)
			metricsMap[NextGC] = Gauge(memStats.NextGC)
			metricsMap[NumForcedGC] = Gauge(memStats.NumForcedGC)
			metricsMap[NumGC] = Gauge(memStats.NumGC)
			metricsMap[OtherSys] = Gauge(memStats.OtherSys)
			metricsMap[PauseTotalNs] = Gauge(memStats.PauseTotalNs)
			metricsMap[StackInuse] = Gauge(memStats.StackInuse)
			metricsMap[StackSys] = Gauge(memStats.StackSys)
			metricsMap[Sys] = Gauge(memStats.Sys)
			metricsMap[PollCount] = Counter(pollCount)
			metricsMap[RandomValue] = Gauge(rand.Intn(10000))
		}
	}(pollInterval)

	go func(host string, port int, interval time.Duration) {
		for {
			<-time.After(interval)
			log.Println("Sending metrics...")

			for key, value := range metricsMap {
				metricName := key

				metricType := "gauge"
				if metricName == PollCount {
					metricType = "counter"
				}

				var metricValue string
				if metricType == "counter" {
					metricValue = fmt.Sprintf("%d", value)
				} else if metricType == "gauge" {
					metricValue = fmt.Sprintf("%f", value)
				}

				url := fmt.Sprintf("http://%s:%d/update/%s/%s/%s",
					host,
					port,
					metricType,
					metricName,
					metricValue,
				)

				applicationType := "text/plain"
				body := []byte(fmt.Sprintf("%f", value))

				request, err := http.Post(url, applicationType, bytes.NewBuffer(body))
				if err != nil {
					log.Fatal(err)
				}

				err = request.Body.Close()
				if err != nil {
					log.Fatal(err)
				}
			}

		}
	}(host, port, reportInterval)

	for {
		time.Sleep(time.Second)
	}
}
