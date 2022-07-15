package main

import (
	"sync"
)

type AgentCache struct {
	metricsMap map[string]*Metric
	mu         *sync.RWMutex
}

func NewAgentCache() *AgentCache {
	return &AgentCache{
		metricsMap: make(map[string]*Metric),
		mu:         &sync.RWMutex{},
	}
}

func (c *AgentCache) Set(key string, value *Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metricsMap[key] = value
}

func (c *AgentCache) Get(key string) (*Metric, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.metricsMap[key]

	return v, ok
}

func (c *AgentCache) MapToSlice() []*Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var metricsSlice []*Metric
	for _, metric := range c.metricsMap {
		metricsSlice = append(metricsSlice, metric)
	}
	return metricsSlice
}
