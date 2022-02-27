package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/handlers"
	"YP-metrics-and-alerting/internal/helpers"
	"YP-metrics-and-alerting/internal/repository"
	"YP-metrics-and-alerting/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	serverAddress := helpers.GetEnv("ADDRESS", "127.0.0.1:8080")
	storeInterval := helpers.StringToSeconds(helpers.GetEnv("STORE_INTERVAL", "0s"))
	storeFile := helpers.GetEnv("STORE_FILE", "devops-metrics-db.json")
	restore, err := strconv.ParseBool(helpers.GetEnv("RESTORE", "true"))
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Config{
		Address:       serverAddress,
		StoreFile:     storeFile,
		StoreInterval: storeInterval,
		Restore:       restore,
	}

	app := &config.Application{
		Config: cfg,
	}

	mapStorage := repository.NewMapStorageRepo()

	repo := handlers.NewRepo(app, mapStorage)

	fileStorage := storage.NewFileStorage(repo.App.StoreFile)
	app.FileStorage = fileStorage

	go repo.ServeFileStorage(fileStorage)

	go handleSignals(repo)

	server := &http.Server{
		Addr:    serverAddress,
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
		repo.Jsonchik()
	default:
		log.Println("unknown signal")
	}

	os.Exit(0)
}
