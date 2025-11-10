package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Danjfreire/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	db             *database.Queries
}

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":8080"
	config := ApiConfig{FileServerHits: atomic.Int32{}, db: dbQueries}

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
