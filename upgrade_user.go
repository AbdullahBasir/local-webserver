package main

import (
	"encoding/json"
	"net/http"

	"github.com/AbdullahBasir/local-webserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) WebhookEvent(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Header, or Header content is empty")
		return
	}
	if apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized Access")
		return
	}

	type eventPayLoad struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	event := eventPayLoad{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&event)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if event.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = cfg.dbQueries.UpgradeUser(r.Context(), event.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Membership resource not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
