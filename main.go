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
	//"github.com/JuanasoKsKs/Chirpy/internal/auth"
	//"fmt"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
	secret string
}




func main() {
	const filepathRoot = "."
	const port = "8080"
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
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
		secret: secret,
	}

	mux.Handle("/app/", cfgs.middlewareMetricsInc(prefixed_hadler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfgs.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfgs.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfgs.handlerChirps)
	mux.HandleFunc("POST /api/users", cfgs.handlerCreateUser)
	mux.HandleFunc("GET /api/chirps", cfgs.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfgs.handlerGetChirp)
	mux.HandleFunc("POST /api/login", cfgs.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfgs.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfgs.handlerRevoke)
	mux.HandleFunc("PUT /api/users", cfgs.handlerUpdate)
	srv := &http.Server{
		Addr : ":" + port,
		Handler : mux,
	}


	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
	
}



