package main

import (
	"YP-metrics-and-alerting/internal/helpers"
	"fmt"
	"sync"
)

type Metric struct {
	Name  string   `json:"id"`
	Type  string   `json:"type"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
	Hash  string   `json:"hash,omitempty"`
	mu    *sync.RWMutex
}

func (m *Metric) SetHash(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var hash string

	if m.Type == GaugeMetricType {
		str := fmt.Sprintf("%s:gauge:%f", m.Name, *m.Value)
		hash = helpers.Hash(str, key)
	} else if m.Type == CounterMetricType {
		str := fmt.Sprintf("%s:counter:%d", m.Name, *m.Delta)
		hash = helpers.Hash(str, key)
	}

	m.Hash = hash
}

func NewGaugeMetric(name string, value float64) *Metric {
	return &Metric{
		Name:  name,
		Type:  GaugeMetricType,
		Value: &value,
		mu:    &sync.RWMutex{},
	}
}

func NewCounterMetric(name string, delta int64) *Metric {
	return &Metric{
		Name:  name,
		Type:  CounterMetricType,
		Delta: &delta,
		mu:    &sync.RWMutex{},
	}
}
