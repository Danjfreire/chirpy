package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Danjfreire/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *ApiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp creation params", err)
		return
	}

	sfwChirpBody := removeProfanity(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: sfwChirpBody, UserID: uuid.MustParse(params.UserId)})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create chirp", err)
		return
	}

	res := Chirp{
		Id:        chirp.ID,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, res)
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
