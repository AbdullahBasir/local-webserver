package main

import (
	"net/http"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) DeleteSingleChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized, unable to get token from Header")
		return
	}

	id := r.PathValue("chirpID")
	parseID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "ID unable to be decoded")
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), parseID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	user, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Token")
		return
	}

	if chirp.UserID != user {
		respondWithError(w, 403, "Invalid authorization to delete")
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), parseID)
	if err != nil {
		respondWithError(w, 500, "Failed to delete chirp")
		return
	}
	w.WriteHeader(204)
}
