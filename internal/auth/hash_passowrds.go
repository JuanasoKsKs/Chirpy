package auth

import(
	"github.com/alexedwards/argon2id"
	"net/http"
	"strings"
	"errors"
)

func HashPassword(password string) (string, error) {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashed, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader:= headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No Authorization header")
	}
	bearer := strings.TrimPrefix(authHeader, "Bearer ")
	return bearer, nil
}