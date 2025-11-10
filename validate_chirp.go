package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type responseVal struct {
		CleanedBody string `json:"cleaned_body"`
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

	respBody := responseVal{CleanedBody: removeProfanity(body.Body)}
	respondWithJSON(w, http.StatusOK, respBody)
}

func removeProfanity(input string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	resultWords := []string{}
	inputWords := strings.Split(input, " ")

	for _, word := range inputWords {
		lowerWord := strings.ToLower(word)
		if _, exists := profaneWords[lowerWord]; exists {
			resultWords = append(resultWords, "****")
		} else {
			resultWords = append(resultWords, word)
		}
	}

	return strings.Join(resultWords, " ")
}
