package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var filter = []string{"kerfuffle", "sharbert", "fornax"}

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanBody := dirtyChirp(params.Body, filter)
	respondWithJSON(w, http.StatusOK, returnVals{CleanBody: cleanBody})
}

func dirtyChirp(body string, filter []string) string {
	cmpStrings := strings.Split(body, " ")
	for i, string := range cmpStrings {
		for _, badWord := range filter {
			if strings.ToLower(string) == badWord {
				cmpStrings[i] = "****"
			}
		}
	}
	result := strings.Join(cmpStrings, " ")
	return result
}
