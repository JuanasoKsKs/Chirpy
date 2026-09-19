package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"time"
	"github.com/JuanasoKsKs/Chirpy/internal/database"
	"github.com/JuanasoKsKs/Chirpy/internal/auth"
	"errors"
	//"fmt"
	//"log"
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
	auth_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No authorization header", err)
		return
	}

	id, err := auth.ValidateJWT(auth_token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Invalid access token", err)
		return
	}


	params.Body = filterProfane(params.Body)
	chirpDB, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: id,
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
func IDFromRequest(r *http.Request) (uuid.UUID, error){
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	userID, err := IDFromRequest(r)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}
	chirpsDB := []database.Chirp{}
	if userID != uuid.Nil {
		chirpsDB, err = cfg.dbQueries.GetChirpsID(r.Context(), userID)
	}else{
		chirpsDB, err = cfg.dbQueries.GetChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err) //500
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
		respondWithError(w, http.StatusInternalServerError, "couldn't parse the UUID string", err) //500
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
func (cfg *apiConfig) handlerDelete(w http.ResponseWriter, r *http.Request) {
	chirpString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpString)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't parse the UUID string", err) // 404
		return
	}
	chirpDB, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't get chirp", err) //404
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Token on authorization header", err) //401
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err) //401
		return
	}
	if chirpDB.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not the owner", errors.New("This user is not the author of the chirp"))
		return
	}
	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpDB.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't Delete the Chirp", err)
	}
	respondWithJSON(w, http.StatusNoContent, response{}) // 204
	
}