package main

import (
	"YP-metrics-and-alerting/internal/config"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestInitAgentConfig(t *testing.T) {
	cfg := InitConfig()
	assert.IsType(t, &config.AgentConfig{}, cfg)
	assert.Equal(t, cfg.Address, "http://127.0.0.1:8080")
	assert.Equal(t, cfg.PoolInterval, 2*time.Second)
	assert.Equal(t, cfg.ReportInterval, 10*time.Second)
}

func TestUpdateMetrics(t *testing.T) {
	interval := time.Second * 1
	wg := &sync.WaitGroup{}
	cache := NewAgentCache()

	wg.Add(1)
	go UpdateMetrics(interval, wg, cache)

	time.Sleep(time.Second * 2)
	allocMetric, ok := cache.Get(Alloc)
	if !ok {
		t.Fail()
	}
	assert.Equal(t, allocMetric.Name, Alloc)
	assert.Equal(t, allocMetric.Type, "gauge")
	assert.NotEqual(t, allocMetric.Delta, nil)
}
