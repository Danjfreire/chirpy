package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Danjfreire/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) upgradeUserToChirpyRedHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetApiKey(r.Header)

	if err != nil || apiKey != cfg.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "invalid API key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	updated, err := cfg.db.UpdateToChirpyRed(r.Context(), uuid.MustParse(params.Data.UserId))
	log.Printf("user is chirpy red: %v\n", updated.IsChirpyRed)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed to upgrade user to Chirpy Red", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
