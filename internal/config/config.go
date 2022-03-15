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
	Key            string
}

type ServerConfig struct {
	Address       string
	StoreFile     string
	Restore       bool
	StoreInterval time.Duration
	Key           string
}

type Application struct {
	Config        ServerConfig
	FileStorage   storage.FileStorage
	TemplateCache map[string]*template.Template
}
