package config

import (
	"YP-metrics-and-alerting/internal/storage"
	"time"
)

type Config struct {
	Address       string
	StoreFile     string
	Restore       bool
	StoreInterval time.Duration
}

type Application struct {
	Config
	FileStorage storage.Storage
}
