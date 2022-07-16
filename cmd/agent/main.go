package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/helpers"
	"flag"
	"fmt"
	"sync"
)

func InitConfig() *config.AgentConfig {
	cfg := &config.AgentConfig{}

	serverHost := helpers.GetEnv("ADDRESS", "127.0.0.1:8080")
	reportInterval := helpers.StringToSeconds(helpers.GetEnv("REPORT_INTERVAL", "10s"))
	pollInterval := helpers.StringToSeconds(helpers.GetEnv("POLL_INTERVAL", "2s"))
	key := helpers.GetEnv("KEY", "")

	flag.StringVar(&cfg.Address, "a", serverHost, "Send metrics to address:port")
	flag.DurationVar(&cfg.ReportInterval, "r", reportInterval, "Report of interval")
	flag.DurationVar(&cfg.PoolInterval, "p", pollInterval, "Pool of interval")
	flag.StringVar(&cfg.Key, "k", key, "Hashing key")
	flag.Parse()

	protocol := "http"
	cfg.Address = fmt.Sprintf("%s://%s", protocol, cfg.Address)

	return cfg
}

func main() {
	cfg := InitConfig()

	cache := NewAgentCache()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go UpdateMetrics(cfg.PoolInterval, wg, cache)

	wg.Add(1)
	go SendMetrics(cfg, wg, cache)

	wg.Wait()
}
