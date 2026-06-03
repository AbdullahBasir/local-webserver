package main

import (
	"net/http"
)

func (cfg *apiConfig) RetrieveChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to retrieve chirps")
		return
	}
	respBody := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		respBody[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}
	respondWithJSON(w, 200, respBody)
}
