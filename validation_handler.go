package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func chirpValidationHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid request body: %v", err))
		return
	} else if len(reqBody.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	wordSplit := strings.Split(reqBody.Body, " ")
	for i, word := range wordSplit {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			wordSplit[i] = "****"
		}
	}
	wordJoin := strings.Join(wordSplit, " ")

	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	valBody := returnVals{
		Cleaned_body: wordJoin,
	}

	err = respondWithJSON(w, 200, valBody)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error marshaling response: %v", err))
		return
	}
}
