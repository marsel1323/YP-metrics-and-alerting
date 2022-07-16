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
	assert.IsType(t, "", gaugeMetric.Hash)

	gaugeMetric.SetHash("hash")
	assert.IsType(t, "", gaugeMetric.Hash)
}

func TestNewCounterMetric(t *testing.T) {
	counterMetric := NewCounterMetric(PollCount, 1.0)
	assert.IsType(t, &Metric{}, counterMetric)
	assert.Equal(t, counterMetric.Name, PollCount)
	assert.Equal(t, *counterMetric.Delta, int64(1))
	assert.IsType(t, "", counterMetric.Hash)

	counterMetric.SetHash("hash")
	assert.IsType(t, "", counterMetric.Hash)
}
