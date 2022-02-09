package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Gauge float64
type Counter int64

type Metrics struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	PollCount     Counter
	RandomValue   Gauge
}

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second
	serverAddress := "localhost:8080"

	var metrics Metrics

	go func(interval time.Duration) {
		var memStats runtime.MemStats
		var pollCount int

		rand.Seed(time.Now().Unix())

		for {
			<-time.After(interval)

			pollCount++

			runtime.ReadMemStats(&memStats)

			metrics = Metrics{
				Alloc:         Gauge(memStats.Alloc),
				BuckHashSys:   Gauge(memStats.BuckHashSys),
				Frees:         Gauge(memStats.Frees),
				GCCPUFraction: Gauge(memStats.GCCPUFraction),
				GCSys:         Gauge(memStats.GCSys),
				HeapAlloc:     Gauge(memStats.HeapAlloc),
				HeapIdle:      Gauge(memStats.HeapIdle),
				HeapInuse:     Gauge(memStats.HeapInuse),
				HeapObjects:   Gauge(memStats.HeapObjects),
				HeapReleased:  Gauge(memStats.HeapReleased),
				HeapSys:       Gauge(memStats.HeapSys),
				LastGC:        Gauge(memStats.LastGC),
				Lookups:       Gauge(memStats.Lookups),
				MCacheInuse:   Gauge(memStats.MCacheInuse),
				MCacheSys:     Gauge(memStats.MCacheSys),
				MSpanInuse:    Gauge(memStats.MSpanInuse),
				MSpanSys:      Gauge(memStats.MSpanSys),
				Mallocs:       Gauge(memStats.Mallocs),
				NextGC:        Gauge(memStats.NextGC),
				NumForcedGC:   Gauge(memStats.NumForcedGC),
				NumGC:         Gauge(memStats.NumGC),
				OtherSys:      Gauge(memStats.OtherSys),
				PauseTotalNs:  Gauge(memStats.PauseTotalNs),
				StackInuse:    Gauge(memStats.StackInuse),
				StackSys:      Gauge(memStats.StackSys),
				Sys:           Gauge(memStats.Sys),
				PollCount:     Counter(pollCount),
				RandomValue:   Gauge(rand.Intn(10000)),
			}
		}
	}(pollInterval)

	go func(serverAddress string, interval time.Duration) {
		host := "127.0.0.1"
		port := 8080

		for {
			<-time.After(interval)

			b, _ := json.Marshal(metrics)

			metricsMap := make(map[string]interface{})

			err := json.Unmarshal(b, &metricsMap)
			if err != nil {
				log.Fatal(err)
			}

			for k, v := range metricsMap {
				metricType := "gauge"
				metricName := k
				if metricName == "PollCount" {
					metricType = "counter"
				}
				metricValue := v
				url := fmt.Sprintf("http://%s:%d/update/%s/%s/%f",
					host,
					port,
					metricType,
					metricName,
					metricValue,
				)
				applicationType := "text/plain"
				body := []byte(fmt.Sprintf("%f", v))
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
	}(serverAddress, reportInterval)

	for {
		time.Sleep(time.Second)
	}
}
