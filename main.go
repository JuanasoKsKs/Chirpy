package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"encoding/json"
)



type apiConfig struct {
	fileserverHits atomic.Int32
}





func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	prefixed_hadler := http.StripPrefix("/app",http.FileServer(http.Dir(filepathRoot)))
	cfgs := &apiConfig{}

	mux.Handle("/app/", cfgs.middlewareMetricsInc(prefixed_hadler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfgs.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfgs.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	srv := &http.Server{
		Addr : ":" + port,
		Handler : mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
	
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	log.Println(params, 2)
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 500, "Couldn't decode parameters", err)
		return
	}

	log.Println(params,3)
	const maxCirpLength = 140

	if len(params.Body) > maxCirpLength {
		responseWithError(w, 401, "Chirp is too long", nil)
		return
	}
	respondWithJSON(w, 200, returnVals{
		Valid: true,
	})


}

func responseWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {

		log.Printf("Responding with 5XX error: %s\n", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-type", "application/son")
    dat, err := json.Marshal(payload)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
	w.WriteHeader(code)
	w.Write(dat)
}