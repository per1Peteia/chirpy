package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/per1Peteia/chirpy/internal/database"
)

type tagChirp struct {
	ID         uuid.UUID `json:"id"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	User_ID    uuid.UUID `json:"user_id"`
}

// validate chirps logic (checks for bad words)
// save chirps if valid logic
var filter = []string{"kerfuffle", "sharbert", "fornax"}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	type returnVals struct {
		tagChirp
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

	// cleanBody will either evaluate to a clean string or a cleaned string
	cleanBody := dirtyChirp(params.Body, filter)
	parsedUserID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing json into uuid")
		return
	}

	chirp, err := cfg.dbQueries.CreateValidChirp(r.Context(), database.CreateValidChirpParams{
		Body:   cleanBody,
		UserID: parsedUserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating chirp record")
		return
	}

	// map the queried struct to a tagged struct
	returnChirp := tagChirp{
		ID:         chirp.ID,
		Created_At: chirp.CreatedAt,
		Updated_At: chirp.UpdatedAt,
		Body:       cleanBody,
		User_ID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, returnVals{returnChirp})
}

// this could be refactored to a map action to be less time complex
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
