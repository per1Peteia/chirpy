package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header set")
	}
	apiKey, _ := strings.CutPrefix(authHeader, "ApiKey ")
	return apiKey, nil
}
