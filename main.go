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
	_ "github.com/lib/pq"
)

type ApiConfig struct {
	FileServerHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	server := &http.Server{}
	server.Handler = mux
	server.Addr = ":8080"
	config := ApiConfig{FileServerHits: atomic.Int32{}, db: dbQueries, platform: platform}

	// app routes
	mux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))

	// api routes
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/users", config.createUserHandler)
	mux.HandleFunc("POST /api/chirps", config.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", config.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpId}", config.findChirpHandler)

	// admin routes
	mux.HandleFunc("GET /admin/metrics", config.metricsHandler)
	mux.HandleFunc("POST /admin/reset", config.resetUsersHandler)

	fmt.Println("Starting server on http://localhost:8080")
	server.ListenAndServe()
}
