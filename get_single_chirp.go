package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetSingleChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	parseID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(r.Context(), parseID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}
	respondWithJSON(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
