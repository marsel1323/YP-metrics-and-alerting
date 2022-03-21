package storage

import (
	"YP-metrics-and-alerting/internal/models"
)

type FileStorage interface {
	Save(data []byte) error
	Retrieve() ([]*models.Metrics, error)
}
