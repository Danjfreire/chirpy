package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Danjfreire/chirpy/internal/auth"
	"github.com/Danjfreire/chirpy/internal/database"
)

type User struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *ApiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid users creation params", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Error", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hashedPassword})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	resp := User{
		Id:          user.ID.String(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func (cfg *ApiConfig) resetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "resetting users is only allowed in dev environment", nil)
		return
	}

	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to reset users", err)
		return
	}
}

func (cfg *ApiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.tokenSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user update params", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Error", err)
		return
	}

	updateUserParams := database.UpdateUserPasswordParams{ID: userId, Email: params.Email, HashedPassword: hashedPassword}
	user, err := cfg.db.UpdateUserPassword(r.Context(), updateUserParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update user", err)
		return
	}

	resp := response{
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
