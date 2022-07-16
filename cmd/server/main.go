package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/gzip"
	"YP-metrics-and-alerting/internal/handlers"
	"YP-metrics-and-alerting/internal/render"
	"YP-metrics-and-alerting/internal/repository"
	"YP-metrics-and-alerting/internal/storage"
	"database/sql"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.InitServerConfig()

	app := &config.Application{
		Config: cfg,
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	render.NewRenderer(app)

	var dbStorage repository.DBRepo
	var db *sql.DB
	if cfg.DSN != "" {
		db, err = initDB(cfg.DSN)
		if err != nil {
			log.Println(err.Error())
		}
		dbStorage = repository.NewPostgresStorage(db)
	} else {
		dbStorage = repository.NewMapStorageRepo()
	}

	repo := handlers.NewRepo(app, dbStorage)

	var fileStorage storage.FileStorage = storage.NewJSONFileStorage(app.Config.StoreFile)
	app.FileStorage = fileStorage
	go repo.ServeFileStorage(fileStorage)

	go handleSignals(repo)

	server := &http.Server{
		Addr:    app.Config.Address,
		Handler: gzip.GzipHandle(Routes(repo)),
	}
	log.Println("Server is serving on", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func initDB(dsn string) (*sql.DB, error) {
	log.Println("Connect to DB:", dsn)
	db, err := sql.Open("pgx", dsn)
	//defer db.Close()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	// create table if not exists
	row := db.QueryRow(`
		CREATE TABLE IF NOT EXISTS metrics
		(
			id    varchar not null,
			type  varchar not null,
			delta bigint,
			value double precision,
			hash  varchar
		);
	`)
	if err := row.Err(); err != nil {
		log.Fatal("Create metrics table error:", err.Error())
	}

	// create index
	row = db.QueryRow(`
		CREATE UNIQUE INDEX IF NOT EXISTS metrics_id_uindex
			ON metrics (id);
	`)
	if err := row.Err(); err != nil {
		log.Fatal("Create metrics table error:", err.Error())
	}

	return db, nil
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
