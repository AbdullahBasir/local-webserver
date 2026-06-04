package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Header not found")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("Header found, but no bearer token found")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}
