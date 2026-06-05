package main

import (
	"net/http"
	"time"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
)

func (cfg *apiConfig) CreateRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Token")
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Token")
		return
	}

	token, err := auth.MakeJWT(user, cfg.Secret, 1*time.Hour)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Token")
		return
	}
	respondWithJSON(w, 200, response{
		Token: token,
	})
}
