package main

import (
	"net/http"
	"time"

	"github.com/gh4rris/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}
	if err = cfg.db.RevokeToken(r.Context(), token); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
