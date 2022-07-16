package main

import (
	"YP-metrics-and-alerting/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func SendMetrics(cfg *config.AgentConfig, wg *sync.WaitGroup, cache *AgentCache) {
	defer wg.Done()

	ticker := time.NewTicker(cfg.ReportInterval)
	for range ticker.C {
		log.Println("Sending metrics...")

		url := fmt.Sprintf("%s/updates", cfg.Address)

		var metricsList = cache.MapToSlice()

		if cfg.Key != "" {
			for _, metric := range metricsList {
				metric.SetHash(cfg.Key)
			}
		}

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
