package main

import (
	"net/http"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) DeleteSingleChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, unable to get token from Header")
		return
	}

	id := r.PathValue("chirpID")
	parseID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID unable to be decoded")
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), parseID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	user, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized Token")
		return
	}

	if chirp.UserID != user {
		respondWithError(w, http.StatusForbidden, "Invalid authorization to delete")
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), parseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
