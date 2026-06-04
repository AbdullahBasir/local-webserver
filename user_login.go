package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password           string `json:"password"`
		Email              string `json:"email"`
		expires_in_seconds int
	}

	type userToken struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
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

	expire := 0
	if params.expires_in_seconds != 0 {
		expire = params.expires_in_seconds * int(time.Second)
	}
	expire = 1 * int(time.Hour)
	jwt, err := auth.MakeJWT(user.ID, cfg.Secret, time.Duration(expire))
	if err != nil {
		respondWithError(w, 500, "Internal Server Error")
	}

	respondWithJSON(w, 200, userToken{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwt,
	})
}
