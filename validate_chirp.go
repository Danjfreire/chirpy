package main

import (
	"encoding/json"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type responseVal struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	body := parameters{}
	err := decoder.Decode(&body)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// failed to parse incoming JSON
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body", err)
	}

	// chirp exceeds 140 characters
	if len(body.Body) > 140 || len(body.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp body must be between 1 and 140 characters", nil)
		return
	}

	respBody := responseVal{Valid: true}
	respondWithJSON(w, http.StatusOK, respBody)
}
