package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

//
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	})
	return newToken.SignedString(signingKey)
}

//
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(t *jwt.Token) (interface{}, error) {return []byte(tokenSecret), nil},
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("Invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid user ID: %w", err)
	}
	return id, nil
}

//
func GetBearerToken(headers http.Header) (string, error) {
	authVal := headers.Get("Authorization")
	if authVal == "" {
		return "", errors.New("Authorization header not found")
	}
	splitAuth := strings.Split(authVal, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("Malformed authorization header")
	}
	return splitAuth[1], nil
}