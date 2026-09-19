package auth

import(
	"github.com/alexedwards/argon2id"
	"net/http"
	"strings"
	"errors"
	"crypto/rand"
	"encoding/hex"
	//"fmt"
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
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("Malformed authorization header")
	}
	return splitAuth[1], nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	encoded := hex.EncodeToString(key)
	return encoded
}