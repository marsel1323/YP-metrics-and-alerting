package config

import (
	"YP-metrics-and-alerting/internal/storage"
	"time"
)

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PoolInterval   time.Duration
}

type ServerConfig struct {
	Address       string
	StoreFile     string
	Restore       bool
	StoreInterval time.Duration
}

type Application struct {
	Config      ServerConfig
	FileStorage storage.Storage
}
