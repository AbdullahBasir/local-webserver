package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
	"github.com/AbdullahBasir/local-webserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized, unable to get token from Header")
		return
	}

	user, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Token")
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request body")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Failed to hash password")
		return
	}

	userInfo, err := cfg.dbQueries.UpdateUserInfo(r.Context(), database.UpdateUserInfoParams{
		Email:    params.Email,
		Password: hashedPassword,
		ID:       user,
	})
	if err != nil {
		respondWithError(w, 500, "Server failed to process request")
		return
	}
	respondWithJSON(w, 200, response{
		ID:        userInfo.ID,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
		Email:     userInfo.Email,
	})
}
