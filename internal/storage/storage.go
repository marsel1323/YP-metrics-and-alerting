package storage

import (
	"YP-metrics-and-alerting/internal/repository"
)

type FileStorage interface {
	Save(data []byte) error
	Retrieve() (*repository.MapStorageRepo, error)
}
