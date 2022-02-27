package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/handlers"
	"YP-metrics-and-alerting/internal/repository"
	"YP-metrics-and-alerting/internal/storage"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.ServerConfig{}

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "Listen to address:port")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "Save metrics to file")
	flag.BoolVar(&cfg.Restore, "r", true, "Restore from file")
	flag.DurationVar(&cfg.StoreInterval, "i", 300*time.Second, "Interval of store to file")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

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
		Addr:    cfg.Address,
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
