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
	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
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

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleaned(reqBody.Body, badWords)

	respondWithJSON(w, 200, returnVals{
		Cleaned_body: cleaned,
	})
}

func getCleaned(body string, badWords map[string]struct{}) string {
	wordSplit := strings.Split(body, " ")
	for i, word := range wordSplit {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			wordSplit[i] = "****"
		}
	}
	cleaned := strings.Join(wordSplit, " ")
	return cleaned
}
