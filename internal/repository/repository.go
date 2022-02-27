package repository

import "errors"

type StorageRepo interface {
	GetAllGaugeMetricValues() (map[string]float64, error)
	GetGaugeMetricValue(string) (float64, error)
	SetGaugeMetricValue(string, float64) error

	GetAllCounterMetricValues() (map[string]int64, error)
	GetCounterMetricValue(string) (int64, error)
	SetCounterMetricValue(string, int64) error
}

type MapStorageRepo struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMapStorageRepo() *MapStorageRepo {
	return &MapStorageRepo{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (m *MapStorageRepo) GetGaugeMetricValue(metricName string) (float64, error) {
	value, ok := m.Gauge[metricName]
	if !ok {
		return 0, errors.New("")
	}
	return value, nil
}

func (m *MapStorageRepo) GetAllGaugeMetricValues() (map[string]float64, error) {
	return m.Gauge, nil
}

func (m *MapStorageRepo) SetGaugeMetricValue(metricName string, metricValue float64) error {
	m.Gauge[metricName] = metricValue
	return nil
}

func (m *MapStorageRepo) GetAllCounterMetricValues() (map[string]int64, error) {
	return m.Counter, nil
}

func (m *MapStorageRepo) GetCounterMetricValue(metricName string) (int64, error) {
	value, ok := m.Counter[metricName]
	if !ok {
		return 0, errors.New("")
	}
	return value, nil
}

func (m *MapStorageRepo) SetCounterMetricValue(metricName string, metricValue int64) error {
	m.Counter[metricName] += metricValue
	return nil
}
