package main

import (
	"encoding/json"
	"net/http"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
)

func (cfg *apiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request body")
		return
	}

	user, err := cfg.dbQueries.UserLogin(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Invalid email or password - Unauthorized")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil || !match {
		respondWithError(w, 401, "Invalid email or password - Unauthorized")
		return
	}

	respondWithJSON(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
