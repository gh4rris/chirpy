package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gh4rris/chirpy/internal/auth"
	"github.com/gh4rris/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err = decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	var badWords = []string{"kerfuffle", "sharbert", "fornax"}
	chirpParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	}
	for _, word := range strings.Split(chirpParams.Body, " ") {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				chirpParams.Body = strings.ReplaceAll(chirpParams.Body, word, "****")
			}
		}
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	var chirps []database.Chirp
	var err error
	if authorID != "" {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID query", err)
			return
		}
		chirps, err = cfg.db.GetUserChirps(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	} else {
		chirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	}
	if sort != "" {
		if sort != "asc" && sort != "desc" {
			respondWithError(w, http.StatusBadRequest, "Invalid sort query", nil)
			return
		}
		if sort == "desc" {
			for i, j := 0, len(chirps)-1; i < j; i, j = i+1, j-1 {
				chirps[i], chirps[j] = chirps[j], chirps[i]
			}
		}
	}
	response := []Chirp{}
	for _, chirp := range chirps {
		item := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		response = append(response, item)
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find Chirp", err)
		return
	}
	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "User not authorized to delete Chirp", nil)
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete Chirp", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
