package main

import (
	"YP-metrics-and-alerting/internal/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	host := "127.0.0.1"
	port := 8080
	serverAddress := fmt.Sprintf("%s:%d", host, port)

	http.HandleFunc("/status", handlers.StatusHandler)
	http.HandleFunc("/update/", handlers.UpdateHandler)

	server := &http.Server{
		Addr: serverAddress,
	}

	log.Println("Server is serving on", server.Addr)
	log.Fatal(server.ListenAndServe())
}
