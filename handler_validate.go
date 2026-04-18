package main

import (
	"log"
	"encoding/json"
	"net/http"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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
		responseWithError(w, 400, "Chirp is too long", nil)
		return
	}
	filtered := filterProfane(params.Body)
	respondWithJSON(w, 200, returnVals{
		CleanedBody: filtered,
	})


}