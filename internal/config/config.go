package config

import (
	"YP-metrics-and-alerting/internal/storage"
	"time"
)

type AgentConfig struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PoolInterval   time.Duration `env:"POOL_INTERVAL" envDefault:"2s"`
}

type ServerConfig struct {
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
}

type Application struct {
	Config      ServerConfig
	FileStorage storage.Storage
}
