package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
	"github.com/AbdullahBasir/local-webserver/internal/database"
	"github.com/google/uuid"
)

func getCleaned(body string, badWords map[string]struct{}) string {
	wordSplit := strings.Split(body, " ")
	for i, word := range wordSplit {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			wordSplit[i] = "****"
		}
	}
	cleaned := strings.Join(wordSplit, " ")
	return cleaned
}

func (cfg *apiConfig) chirpWriter(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid request body: %v", err))
		return
	} else if len(reqBody.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleaned(reqBody.Body, badWords)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Token")
		return
	}

	create_chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to create chirp: %v", err))
		return
	}
	respondWithJSON(w, 201, Chirp{
		ID:        create_chirp.ID,
		CreatedAt: create_chirp.CreatedAt,
		UpdatedAt: create_chirp.UpdatedAt,
		Body:      create_chirp.Body,
		UserID:    create_chirp.UserID,
	})
}
