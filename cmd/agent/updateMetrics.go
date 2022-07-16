package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func UpdateMetrics(interval time.Duration, wg *sync.WaitGroup, cache *AgentCache) {
	defer wg.Done()

	var pollCount int

	rand.Seed(time.Now().Unix())

	ticker := time.NewTicker(interval)
	for range ticker.C {
		log.Println("Update metrics...")
		pollCount++
		cache.Set(PollCount, NewCounterMetric(PollCount, int64(pollCount)))
		cache.Set(RandomValue, NewGaugeMetric(RandomValue, float64(rand.Intn(10000))))
		SetRuntimeStats(cache)

		log.Println("UpdateExtraMetrics...")
		SetVirtualMemoryStats(cache)
		SetCPUStats(cache)
	}
}
