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

var (
	serverAddressFlag = flag.String("a", "127.0.0.1:8080", "Listen to address:port")
	storeIntervalFlag = flag.String("i", "300s", "Interval of store to file")
	storeFileFlag     = flag.String("f", "/tmp/devops-metrics-db.json", "Save metrics to file")
	restoreFlag       = flag.String("r", "true", "Restore from file")
)

func main() {
	log.Println(os.Args)

	flag.Parse()

	serverAddress := helpers.GetEnv("ADDRESS", *serverAddressFlag)
	storeInterval := helpers.StringToSeconds(helpers.GetEnv("STORE_INTERVAL", *storeIntervalFlag))
	storeFile := helpers.GetEnv("STORE_FILE", *storeFileFlag)
	restore, err := strconv.ParseBool(helpers.GetEnv("RESTORE", *restoreFlag))
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.ServerConfig{
		Address:       serverAddress,
		StoreFile:     storeFile,
		Restore:       restore,
		StoreInterval: storeInterval,
	}

	log.Println(cfg)
	log.Println("restore", cfg.Restore)

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
