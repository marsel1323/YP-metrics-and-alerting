package repository

import "YP-metrics-and-alerting/internal/models"

type DBRepo interface {
	GetMetric(id string) (*models.Metrics, error)
	SetMetric(metric *models.Metrics) error
	GetMetricsList() ([]*models.Metrics, error)
	SetMetricsList(metricsList []*models.Metrics) error
	SetMetricsListFromFile(metricsList []*models.Metrics) error
}
