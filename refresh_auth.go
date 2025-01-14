package main

import (
	"net/http"
	"time"

	"github.com/per1Peteia/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	rTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error getting token string")
		return
	}

	rToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), rTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token not found/expired/revoked")
		return
	}

	aToken, err := auth.MakeJWT(rToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error making JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: aToken})
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	rTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error getting token string")
		return
	}

	_, err = cfg.dbQueries.RevokeRefreshToken(r.Context(), rTokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating refresh token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
