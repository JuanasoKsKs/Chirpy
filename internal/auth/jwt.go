package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/google/uuid"
	"errors"
	"fmt"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	method := jwt.SigningMethodHS256
	token := jwt.NewWithClaims(method, jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func (t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	IDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("wrong issuer")
	}
	id, err := uuid.Parse(IDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("wrong user ID: %w", err)
	}
	return id, nil
}

