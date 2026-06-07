package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Header not found")
	}

	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", errors.New("Header found, but no bearer token found")
	}

	key := strings.TrimPrefix(authHeader, "ApiKey ")
	return key, nil
}
