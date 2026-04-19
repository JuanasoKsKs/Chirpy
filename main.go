package main

import (
	"net/http"
	"log"
	"sync/atomic"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"os"
	"database/sql"
	"github.com/JuanasoKsKs/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
}



func main() {
	const filepathRoot = "."
	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening conection to the database: %s\n", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	prefixed_hadler := http.StripPrefix("/app",http.FileServer(http.Dir(filepathRoot)))
	cfgs := &apiConfig{
		dbQueries: database.New(db),
		platform: platform,
	}

	mux.Handle("/app/", cfgs.middlewareMetricsInc(prefixed_hadler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfgs.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfgs.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfgs.handlerChirps)
	mux.HandleFunc("POST /api/users", cfgs.handlerCreateUser)
	mux.HandleFunc("GET /api/chirps", cfgs.handlerGetChirps)
	srv := &http.Server{
		Addr : ":" + port,
		Handler : mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
	
}



