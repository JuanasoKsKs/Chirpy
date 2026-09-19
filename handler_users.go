package main

import(
	"net/http"
	"time"
	"github.com/google/uuid"
	"encoding/json"
	"github.com/JuanasoKsKs/Chirpy/internal/auth"
	"github.com/JuanasoKsKs/Chirpy/internal/database"
	"errors"
)

type User struct {
	ID 					uuid.UUID 	`json:"id"`
	CreatedAt 			time.Time 	`json:"created_at"`
	UpdatedAt	 		time.Time 	`json:"updated_at"`
	Email 				string 		`json:"email"`
	Password			string 		`json:"-"`
	IsChirpyRed			bool		`json:"is_chirpy_red"`
}
type parameters struct {
	Password 	string `json:"password"`
	Email 		string `json:`
}
type response struct {
	User
	Token 			string `json:"token"`
	RefreshToken	string `json:"refresh_token"`
}
type webhooksParams struct {
	Data struct{
		UserID	uuid.UUID	`json:"user_id"`
	}	`json:"data"`
	Event string			`json:"event"`	
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't hashed password", err)
		return
	}
	userDB, err:= cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating the User", err)
		return
	}	
	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID: userDB.ID,
			CreatedAt: userDB.CreatedAt,
			UpdatedAt: userDB.UpdatedAt,
			Email: userDB.Email,
			IsChirpyRed: userDB.IsChirpyRed,
		},
	})
}
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't Decode request", err)
		return
	}
	userDB, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't Get User by Email", err)
		return
	}
	passed, err := auth.CheckPasswordHash(params.Password, userDB.HashedPassword)
	if err != nil || !passed {
		respondWithError(w, 401, "Error comparing hash and password", err)
		return
	}
	accessToken, err := auth.MakeJWT(userDB.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Error generating the jwt", err)
		return
	}
	refreshToken := auth.MakeRefreshToken()
	err = cfg.dbQueries.CreateToken(r.Context(), database.CreateTokenParams{
		Token: refreshToken,
		UserID: userDB.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Error Saving the Token in database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:			userDB.ID,
			UpdatedAt:	userDB.UpdatedAt,
			CreatedAt: 	userDB.CreatedAt,
			Email: 		userDB.Email,
			IsChirpyRed: userDB.IsChirpyRed,
		},
		Token: 			accessToken,
		RefreshToken: 	refreshToken,
	})
}
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No Token on authorization header", err)
		return
	}
	refreshToken, err := cfg.dbQueries.GetToken(r.Context(), token)
	if err != nil || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token invalid or revoked", err)
		return
	}
	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error generating the jwt", err) //401
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token:	accessToken,
	})
}
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No Token on authorization header", err) //400
		return
	}
	err = cfg.dbQueries.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 500, "Error removing token", err)
	}
	respondWithJSON(w, http.StatusNoContent, response{}) // 204
}
func (cfg *apiConfig) handlerUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Token on authorization header", err) //401
		return
	}
	UserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err) //401
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	newHashedPassword, err := auth.HashPassword(params.Password)
	userDB, err := cfg.dbQueries.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		Email: params.Email,
		HashedPassword: newHashedPassword,
		ID: UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't insert credentials into database", err)
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:			userDB.ID,
			UpdatedAt:	userDB.UpdatedAt,
			CreatedAt: 	userDB.CreatedAt,
			Email: 		userDB.Email,
			IsChirpyRed: userDB.IsChirpyRed,
		},
	})

}
func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No key on authorization header", err)
		return
	}
	if key != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Keys don't match", errors.New("Unauthorized user"))
	}
	decoder := json.NewDecoder(r.Body)
	params := webhooksParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "No valid request", errors.New("No valid event - Payment Failed"))
		return
	}
	err = cfg.dbQueries.UpdateUserRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, response{})

}