package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/handlers"
	"YP-metrics-and-alerting/internal/helpers"
	"YP-metrics-and-alerting/internal/render"
	"YP-metrics-and-alerting/internal/repository"
	"YP-metrics-and-alerting/internal/storage"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	serverAddressFlag := flag.String("a", "127.0.0.1:8080", "Listen to address:port")
	storeIntervalFlag := flag.String("i", "300s", "Interval of store to file")
	storeFileFlag := flag.String("f", "/tmp/devops-metrics-db.json", "Save metrics to file")
	restoreFlag := flag.String("r", "true", "Restore from file")
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

	app := &config.Application{
		Config: cfg,
	}

	t, err := template.ParseFiles("./internal/templates/metrics.gohtml")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(t)

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return
	}

	app.TemplateCache = tc

	mapStorage := repository.NewMapStorageRepo()
	fileStorage := storage.NewFileStorage(app.Config.StoreFile)

	repo := handlers.NewRepo(app, mapStorage)
	render.NewRenderer(app)

	app.FileStorage = fileStorage

	go repo.ServeFileStorage(fileStorage)

	go handleSignals(repo)

	server := &http.Server{
		Addr:    app.Config.Address,
		Handler: handlers.GzipHandle(Routes(repo)),
	}
	log.Println("Server is serving on", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func handleSignals(repo *handlers.Repository) {
	captureSignal := make(chan os.Signal, 1)
	signal.Notify(captureSignal, syscall.SIGINT, syscall.SIGTERM)
	time.Sleep(1 * time.Second)

	switch <-captureSignal {
	case syscall.SIGINT:
	case syscall.SIGTERM:
		repo.SaveMetrics()
	default:
		log.Println("unknown signal")
	}

	os.Exit(0)
}
