package main

import (
	"net/http"

	"github.com/Joad/rss_aggregator/internal/auth"
	"github.com/Joad/rss_aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKeyFromHeader(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "API key not found")
			return
		}
		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}
		handler(w, r, user)
	}
}
