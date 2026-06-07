package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) RetrieveChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")
	parseID, err := uuid.Parse(author)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	if author != "" {
		authorChirp, err := cfg.dbQueries.GetAuthorChirps(r.Context(), parseID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps")
			return
		}
		respBody := make([]Chirp, len(authorChirp))
		for i, author := range authorChirp {
			respBody[i] = Chirp{
				ID:        author.ID,
				CreatedAt: author.CreatedAt,
				UpdatedAt: author.UpdatedAt,
				Body:      author.Body,
				UserID:    author.UserID,
			}
		}
		respondWithJSON(w, http.StatusOK, respBody)
		return
	}

	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps")
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
	respondWithJSON(w, http.StatusOK, respBody)
}
