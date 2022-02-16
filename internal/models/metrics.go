package models

type Gauge float64
type Counter int64

type GaugeMetric struct {
	Value float64
}

type CounterMetric struct {
	Value int64
}
