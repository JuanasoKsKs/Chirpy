package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"time"
	"github.com/JuanasoKsKs/Chirpy/internal/database"
)
type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	const maxCirpLength = 140
	if len(params.Body) > maxCirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	params.Body = filterProfane(params.Body)
	chirpDB, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body: chirpDB.Body,
		UserID: chirpDB.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpsDB, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}
	chirps := []Chirp{}
	for i, c := range chirpsDB {
		chirps = append(chirps, Chirp{})
		chirps[i].ID = c.ID
		chirps[i].CreatedAt = c.CreatedAt
		chirps[i].UpdatedAt = c.UpdatedAt
		chirps[i].Body = c.Body
		chirps[i].UserID = c.UserID
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse the UUID string", err)
		return
	}
	chirpDB, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't get chirp", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID: chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body: chirpDB.Body,
		UserID: chirpDB.UserID,
	})
}