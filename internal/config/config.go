package config

import (
	"YP-metrics-and-alerting/internal/storage"
	"html/template"
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
	Config        ServerConfig
	FileStorage   storage.Storage
	TemplateCache map[string]*template.Template
}
