package auth

import (
	"errors"
	"net/http"
	"strings"
)

//
func GetAPIKey(headers http.Header) (string, error) {
	authVal := headers.Get("Authorization")
	if authVal == "" {
		return "", errors.New("Authorization header not found")
	}
	splitAuth := strings.Split(authVal, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("Malformed authorization header")
	}
	apiKey := splitAuth[1]
	return apiKey, nil
}