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

type mapStorageRepo struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMapStorageRepo() *mapStorageRepo {
	return &mapStorageRepo{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (m *mapStorageRepo) GetGaugeMetricValue(metricName string) (float64, error) {
	value, ok := m.Gauge[metricName]
	if !ok {
		return 0, errors.New("")
	}
	return value, nil
}

func (m *mapStorageRepo) GetAllGaugeMetricValues() (map[string]float64, error) {
	return m.Gauge, nil
}

func (m *mapStorageRepo) SetGaugeMetricValue(metricName string, metricValue float64) error {
	m.Gauge[metricName] = metricValue
	return nil
}

func (m *mapStorageRepo) GetAllCounterMetricValues() (map[string]int64, error) {
	return m.Counter, nil
}

func (m *mapStorageRepo) GetCounterMetricValue(metricName string) (int64, error) {
	value, ok := m.Counter[metricName]
	if !ok {
		return 0, errors.New("")
	}
	return value, nil
}

func (m *mapStorageRepo) SetCounterMetricValue(metricName string, metricValue int64) error {
	m.Counter[metricName] += metricValue
	return nil
}
