package main

import(
	"net/http"
	"time"
	"github.com/google/uuid"
	"encoding/json"
)

type User struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type Email struct {
		
	}
	decoder := json.NewDecoder(r.Body)
	user := User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	userDB, err:= cfg.dbQueries.CreateUser(r.Context(), user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating the User", err)
		return
	}
	respondWithJSON(w, 201, User{
		Id: userDB.ID,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Email: userDB.Email,
	})

}

