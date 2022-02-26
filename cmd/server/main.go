package main

import (
	"YP-metrics-and-alerting/internal/config"
	"YP-metrics-and-alerting/internal/helpers"
	"log"
	"net/http"
)

func main() {
	serverAddress := helpers.GetEnv("ADDRESS", "127.0.0.1:8080")

	var cfg config.Config

	cfg.Host = "127.0.0.1"
	cfg.Port = 8080

	//serverAddress := fmt.Sprintf("%s:%d",
	//	cfg.Host,
	//	cfg.Port,
	//)

	app := &config.Application{
		Config: cfg,
	}

	server := &http.Server{
		Addr:    serverAddress,
		Handler: Routes(app),
	}

	log.Println("Server is serving on", server.Addr)
	log.Fatal(server.ListenAndServe())
}
