package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, status int, msg string) {
	type response struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, status, response{
		Error: msg,
	})
}
