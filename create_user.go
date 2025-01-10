package main

import (
	"encoding/json"
	"net/http"

	"github.com/per1Peteia/chirpy/internal/database"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding JSON")
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)

}
