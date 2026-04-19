package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden action, needs to be dev to perform", nil)
		return
	}
	type infoResponse struct {
		Message string `json:"message"`
	}
	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Reseting users table", err)
	}
	respondWithJSON(w, 200, infoResponse{
		Message: "user database was reseted",
	})
}
