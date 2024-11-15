package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp := parameters{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	var badWords = []string{"kerfuffle", "sharbert", "fornax"}
	response := returnVals{
		CleanedBody: chirp.Body,
	}
	for _, word := range strings.Split(response.CleanedBody, " ") {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				response.CleanedBody = strings.ReplaceAll(response.CleanedBody, word, "****")
			}
		}
	}
	respondWithJSON(w, http.StatusOK, response)
}
