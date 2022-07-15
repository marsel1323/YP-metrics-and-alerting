package config

import (
	"YP-metrics-and-alerting/internal/helpers"
	"YP-metrics-and-alerting/internal/storage"
	"flag"
	"html/template"
	"log"
	"strconv"
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
	DSN           string
}

func InitServerConfig() ServerConfig {
	serverAddressFlag := flag.String("a", "127.0.0.1:8080", "Listen to address:port")
	storeIntervalFlag := flag.String("i", "300s", "Interval of store to file")
	storeFileFlag := flag.String("f", "/tmp/devops-metrics-db.json", "Save metrics to file")
	restoreFlag := flag.String("r", "true", "Restore from file")
	keyFlag := flag.String("k", "", "Hashing key")
	dbDsnFlag := flag.String("d", "", "Database DSN")
	flag.Parse()

	serverAddress := helpers.GetEnv("ADDRESS", *serverAddressFlag)
	storeInterval := helpers.StringToSeconds(helpers.GetEnv("STORE_INTERVAL", *storeIntervalFlag))
	storeFile := helpers.GetEnv("STORE_FILE", *storeFileFlag)
	restore, err := strconv.ParseBool(helpers.GetEnv("RESTORE", *restoreFlag))
	if err != nil {
		log.Fatal(err)
	}
	key := helpers.GetEnv("KEY", *keyFlag)
	dbDsn := helpers.GetEnv("DATABASE_DSN", *dbDsnFlag)

	cfg := ServerConfig{
		Address:       serverAddress,
		StoreFile:     storeFile,
		Restore:       restore,
		StoreInterval: storeInterval,
		Key:           key,
		DSN:           dbDsn,
	}

	return cfg
}

type Application struct {
	Config        ServerConfig
	FileStorage   storage.FileStorage
	TemplateCache map[string]*template.Template
}
