package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Danjfreire/chirpy/internal/auth"
)

func (cfg *ApiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"` // in seconds
	}

	type response struct {
		Token string `json:"token"`
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Expected login and password", err)
		return
	}

	user, err := cfg.db.FindUserByEmail(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)

	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Invalid login or password", nil)
	}

	expiration := int32(3600) // default 1 hour

	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds <= 3600 {
		expiration = int32(params.ExpiresInSeconds)
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Second*time.Duration(expiration))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token", err)
		return
	}

	resp := response{
		Token: token,
		User: User{
			Id:        user.ID.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	}

	respondWithJSON(w, http.StatusOK, resp)
}
