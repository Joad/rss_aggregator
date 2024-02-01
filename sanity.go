package main

import "net/http"

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Status: "ok",
	})
}

func errHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
