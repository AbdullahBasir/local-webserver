package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
	"github.com/AbdullahBasir/local-webserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type responseBody struct {
		Id           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
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

	exp := 60 * (24 * time.Hour)
	refreshToken := auth.MakeRefreshToken()

	rToken, err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(exp),
	})
	if err != nil {
		respondWithError(w, 401, "Unauthorized Token")
		return
	}

	expire := 1 * int(time.Hour)
	jwt, err := auth.MakeJWT(user.ID, cfg.Secret, time.Duration(expire))
	if err != nil {
		respondWithError(w, 500, "Internal Server Error")
		return
	}

	respondWithJSON(w, 200, responseBody{
		Id:           rToken.UserID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        jwt,
		RefreshToken: rToken.Token,
	})
}
