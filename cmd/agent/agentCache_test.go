package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestAgentCache(t *testing.T) {
	agentCache := &AgentCache{
		metricsMap: make(map[string]*Metric),
		mu:         &sync.RWMutex{},
	}

	gaugeMetric := NewGaugeMetric("test", 1.0)

	agentCache.Set("test", gaugeMetric)

	received, ok := agentCache.Get("test")
	if !ok {
		t.Fail()
	}

	assert.Equal(t, gaugeMetric.Value, received.Value)
	assert.Equal(t, gaugeMetric.Name, received.Name)
}

func TestNewAgentCache(t *testing.T) {
	agentCache := NewAgentCache()

	gaugeMetric := NewGaugeMetric("test", 1.0)

	agentCache.Set("test", gaugeMetric)

	received, ok := agentCache.Get("test")
	if !ok {
		t.Fail()
	}

	assert.Equal(t, gaugeMetric.Value, received.Value)
	assert.Equal(t, gaugeMetric.Name, received.Name)
}

func TestGaugeMetric(t *testing.T) {
	agentCache := NewAgentCache()

	alloc := NewGaugeMetric(Alloc, 1.0)
	agentCache.Set(Alloc, alloc)

	received, ok := agentCache.Get(alloc.Name)
	if !ok {
		t.Fail()
	}

	assert.Equal(t, alloc.Value, received.Value)
	assert.Equal(t, alloc.Name, received.Name)
}

func TestCounterMetric(t *testing.T) {
	agentCache := NewAgentCache()

	pollCount := NewCounterMetric(PollCount, 1.0)
	agentCache.Set(PollCount, pollCount)

	received, ok := agentCache.Get(pollCount.Name)
	if !ok {
		t.Fail()
	}

	assert.Equal(t, pollCount.Value, received.Value)
	assert.Equal(t, pollCount.Name, received.Name)
}

func TestMapToSlice(t *testing.T) {
	agentCache := NewAgentCache()

	alloc := NewGaugeMetric(Alloc, 1.0)
	pollCount := NewCounterMetric(PollCount, 1.0)

	agentCache.Set(Alloc, alloc)
	agentCache.Set(PollCount, pollCount)

	var want []*Metric
	want = make([]*Metric, 0)
	want = append(want, alloc)
	want = append(want, pollCount)

	metricsSlice := agentCache.MapToSlice()
	assert.ElementsMatch(t, metricsSlice, want)
}

func BenchmarkAgentCache(b *testing.B) {
	agentCache := NewAgentCache()
	for i := 0; i < b.N; i++ {
		testMetric := NewGaugeMetric("test", 1.0)

		agentCache.Set("test", testMetric)

		received, ok := agentCache.Get("test")
		if !ok {
			b.Fail()
		}

		assert.Equal(b, testMetric.Value, received.Value)
		assert.Equal(b, testMetric.Name, received.Name)
	}
}
