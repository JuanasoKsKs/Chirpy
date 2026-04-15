package main

import (
	"net/http"
	"log"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	message := []byte("OK")
	w.Write(message)
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app",http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", testHandler)
	srv := &http.Server{
		Addr : ":" + port,
		Handler : mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
	
}