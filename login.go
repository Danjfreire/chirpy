package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Danjfreire/chirpy/internal/auth"
	"github.com/Danjfreire/chirpy/internal/database"
)

func (cfg *ApiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Second*time.Duration(expiration))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token", err)
		return
	}

	refreshtokenDuration := time.Hour * 24 * 60 // 60 days
	refreshTokenParams := database.CreateRefreshTokenParams{Token: auth.MakeRefreshToken(), UserID: user.ID, ExpiresAt: time.Now().Add(refreshtokenDuration)}
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)

	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Could not create refresh token", err)
		return
	}

	resp := response{
		Token:        token,
		RefreshToken: refreshToken,
		User: User{
			Id:          user.ID.String(),
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *ApiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid refresh token", err)
	}

	refreshToken, err := cfg.db.GetRefreshTokenByToken(r.Context(), token)

	if err != nil || refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token", err)
	}

	expiration := int32(3600) // default 1 hour
	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret, time.Duration(expiration)*time.Second)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
	}

	resp := response{
		Token: newToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *ApiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid refresh token", err)
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
