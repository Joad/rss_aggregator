package main

import (
	"net/http"
	"strconv"

	"github.com/Joad/rss_aggregator/internal/database"
)

func (cfg *apiConfig) getPostsByUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	var err error

	limit := 5
	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "invalid limit")
			return
		}
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting posts")
		return
	}

	respondWithJSON(w, http.StatusOK, dbPostsToPosts(posts))
}
