package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gh4rris/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: accessToken,
	})
}
