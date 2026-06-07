package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
	"github.com/AbdullahBasir/local-webserver/internal/database"
)

func (cfg *apiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to hash password: %v", err))
		return
	}

	create_user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:    params.Email,
		Password: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	respondWithJSON(w, 201, User{
		ID:          create_user.ID,
		CreatedAt:   create_user.CreatedAt,
		UpdatedAt:   create_user.UpdatedAt,
		Email:       create_user.Email,
		IsChirpyRed: create_user.IsChirpyRed,
	})

}
