package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Danjfreire/chirpy/internal/auth"
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
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp creation params", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.tokenSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	sfwChirpBody := removeProfanity(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: sfwChirpBody, UserID: userId})

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

func (cfg *ApiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get chirps", err)
		return
	}

	res := []Chirp{}

	for _, chirp := range chirps {
		res = append(res, Chirp{
			Id:        chirp.ID,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) findChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpId")
	chirp, err := cfg.db.GetChirpByID(r.Context(), uuid.MustParse(chirpId))

	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", nil)
		return
	}

	res := Chirp{
		Id:        chirp.ID,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, res)
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
