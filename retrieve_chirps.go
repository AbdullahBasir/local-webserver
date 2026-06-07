package main

import (
	"net/http"
	"sort"

	"github.com/AbdullahBasir/local-webserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) RetrieveChirps(w http.ResponseWriter, r *http.Request) {
	var parseID uuid.UUID
	var err error

	author := r.URL.Query().Get("author_id")
	if author != "" {
		parseID, err = uuid.Parse(author)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
			return
		}
	}

	var dbChirps []database.Chirp

	if parseID != uuid.Nil {
		dbChirps, err = cfg.dbQueries.GetAuthorChirps(r.Context(), parseID)
	} else {
		dbChirps, err = cfg.dbQueries.GetChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps")
		return
	}
	respBody := []Chirp{}
	for _, dbchirp := range dbChirps {
		respBody = append(respBody, Chirp{
			ID:        dbchirp.ID,
			CreatedAt: dbchirp.CreatedAt,
			UpdatedAt: dbchirp.UpdatedAt,
			Body:      dbchirp.Body,
			UserID:    dbchirp.UserID,
		})
	}

	sorting := r.URL.Query().Get("sort")

	if sorting == "desc" {
		sort.Slice(respBody, func(i, j int) bool {
			return respBody[i].CreatedAt.After(respBody[j].CreatedAt)
		})
		respondWithJSON(w, http.StatusOK, respBody)
		return
	} else {
		sort.Slice(respBody, func(i, j int) bool {
			return respBody[i].CreatedAt.Before(respBody[j].CreatedAt)
		})
		respondWithJSON(w, http.StatusOK, respBody)
	}
}
