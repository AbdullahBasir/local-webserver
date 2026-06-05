package main

import (
	"net/http"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
)

func (cfg *apiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized, could not obtain token from Header")
		return
	}
	err = cfg.dbQueries.RevokeDatabaseToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorized, Missing credentials")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
