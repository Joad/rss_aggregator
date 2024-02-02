package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Joad/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding json")
		return
	}

	now := time.Now().UTC()
	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, dbFeedToFeed(feed))
}
