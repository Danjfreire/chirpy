package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	splits := strings.Split(authHeader, "ApiKey ")

	if len(splits) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	return splits[1], nil
}
