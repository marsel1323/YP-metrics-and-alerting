package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/handlers"
	"YP-metrics-and-alerting/internal/helpers"
	"YP-metrics-and-alerting/internal/repository"
	"YP-metrics-and-alerting/internal/storage"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	cfg := config.ServerConfig{}

	serverAddress := helpers.GetEnv("ADDRESS", "127.0.0.1:8080")
	storeInterval := helpers.StringToSeconds(helpers.GetEnv("STORE_INTERVAL", "30s"))
	storeFile := helpers.GetEnv("STORE_FILE", "/tmp/devops-metric-db.json")
	restore, err := strconv.ParseBool(helpers.GetEnv("RESTORE", "true"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("restore", restore)

	flag.StringVar(&cfg.Address, "a", serverAddress, "Listen to address:port")
	flag.DurationVar(&cfg.StoreInterval, "i", storeInterval, "Interval of store to file")
	flag.StringVar(&cfg.StoreFile, "f", storeFile, "Save metrics to file")
	flag.BoolVar(&cfg.Restore, "r", restore, "Restore from file")
	flag.Parse()
	log.Println(cfg)
	log.Println("restore", &cfg.Restore)
	app := &config.Application{
		Config: cfg,
	}

	mapStorage := repository.NewMapStorageRepo()

	repo := handlers.NewRepo(app, mapStorage)

	fileStorage := storage.NewFileStorage(repo.App.Config.StoreFile)
	app.FileStorage = fileStorage

	go repo.ServeFileStorage(fileStorage)

	go handleSignals(repo)

	server := &http.Server{
		Addr:    app.Config.Address,
		Handler: Routes(repo),
	}
	log.Println("Server is serving on", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func handleSignals(repo *handlers.Repository) {
	var captureSignal = make(chan os.Signal, 1)
	signal.Notify(captureSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	time.Sleep(1 * time.Second)

	switch <-captureSignal {
	case syscall.SIGINT:
		repo.SaveMetrics()
	default:
		log.Println("unknown signal")
	}

	os.Exit(0)
}
