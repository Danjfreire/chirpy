package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":8080"
	config := ApiConfig{FileServerHits: atomic.Int32{}}

	// app routes
	mux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))

	// api routes
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	// admin routes
	mux.HandleFunc("GET /admin/metrics", config.metricsHandler)
	mux.HandleFunc("POST /admin/reset", config.resetMetricsHandler)

	fmt.Println("Starting server on http://localhost:8080")
	server.ListenAndServe()
}
