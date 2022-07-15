package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGaugeMetric(t *testing.T) {
	gaugeMetric := NewGaugeMetric(Alloc, 1.0)

	assert.IsType(t, &Metric{}, gaugeMetric)
	assert.Equal(t, gaugeMetric.Name, Alloc)
	assert.Equal(t, gaugeMetric.Type, "gauge")
	assert.Equal(t, *gaugeMetric.Value, 1.0)
	assert.Equal(t, gaugeMetric.Delta, nil)
	assert.IsType(t, "", gaugeMetric.Hash)
	t.Log(gaugeMetric)

	gaugeMetric.SetHash("hash")
	t.Log(gaugeMetric)

	assert.IsType(t, "", gaugeMetric.Hash)
}

func TestNewCounterMetric(t *testing.T) {
	counterMetric := NewCounterMetric(PollCount, 1.0)
	assert.IsType(t, &Metric{}, counterMetric)
	assert.Equal(t, counterMetric.Name, Alloc)
	assert.Equal(t, *counterMetric.Value, 1.0)
	assert.IsType(t, "", counterMetric.Hash)
	t.Log(counterMetric)
}
