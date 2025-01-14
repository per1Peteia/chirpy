package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/per1Peteia/chirpy/internal/auth"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	parsedID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing path parameter")
		return
	}

	// authentication
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authentication header set")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token signature")
		return
	}

	// authorization
	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error getting chirp record")
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "missing authorization")
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting chirp record")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
