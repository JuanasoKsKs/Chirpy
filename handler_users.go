package main

import(
	"net/http"
	"time"
	"github.com/google/uuid"
	"encoding/json"
	"github.com/JuanasoKsKs/Chirpy/internal/auth"
	"github.com/JuanasoKsKs/Chirpy/internal/database"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't hashed password", err)
	}
	userDB, err:= cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: user.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating the User", err)
		return
	}	
	respondWithJSON(w, http.StatusCreated, User{
		ID: userDB.ID,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Email: userDB.Email,
	})

}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't Decode request", err)
		return
	}
	userDB, err := cfg.dbQueries.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't Get User by Email", err)
		return
	}
	passed, err := auth.CheckPasswordHash(user.Password, userDB.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Error comparing hash and password", err)
		return
	}
	if passed {
		respondWithJSON(w, http.StatusOK, User{
			ID: userDB.ID,
			UpdatedAt: userDB.UpdatedAt,
			CreatedAt: userDB.CreatedAt,
			Email: userDB.Email,
		})
	} else {
		respondWithError(w, 401, "Incorrect Email or password", nil)
		return
	}
	
}
