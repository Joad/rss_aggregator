package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrApiKeyMissing = errors.New("API key not found in header")

func GetApiKeyFromHeader(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	apiKey, found := strings.CutPrefix(authHeader, "ApiKey ")
	if !found {
		return "", ErrApiKeyMissing
	}
	return apiKey, nil
}
