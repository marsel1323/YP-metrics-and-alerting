package repository

import (
	"YP-metrics-and-alerting/internal/models"
	"errors"
	"fmt"
)

type MapStorageRepo map[string]*models.Metrics

func NewMapStorageRepo() MapStorageRepo {
	return make(map[string]*models.Metrics)
}

func (m MapStorageRepo) GetMetric(id string) (*models.Metrics, error) {
	metric, ok := m[id]
	if !ok {
		return nil, fmt.Errorf("metric '%s' not found", id)
	}
	return metric, nil
}

func (m MapStorageRepo) GetMetricsList() ([]*models.Metrics, error) {
	var metrics []*models.Metrics
	for _, metric := range m {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (m MapStorageRepo) SetMetric(metric *models.Metrics) error {
	foundMetric, ok := m[metric.ID]
	if ok {
		if metric.MType != foundMetric.MType {
			return errors.New("metrics types don't match")
		}
		if metric.MType == models.CounterType {
			sum := *foundMetric.Delta + *metric.Delta
			foundMetric.Delta = &sum
		} else if metric.MType == models.GaugeType {
			foundMetric.Value = metric.Value
		}
	} else {
		m[metric.ID] = metric
	}

	return nil
}

func (m MapStorageRepo) SetMetricsList(metricsList []*models.Metrics) error {
	for _, metric := range metricsList {
		m[metric.ID] = metric
	}
	return nil
}

func (m MapStorageRepo) SetMetricsListFromFile(metricsList []*models.Metrics) error {
	for _, metric := range metricsList {
		m[metric.ID] = metric
	}
	return nil
}
