package models

type GaugeMetric struct {
	Value float64
}

type CounterMetric struct {
	Value int64
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
