package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/per1Peteia/chirpy/internal/auth"
	"github.com/per1Peteia/chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		taggedUser
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding json")
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if err = auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expirationTime := time.Hour

	accessTokenString, err := auth.MakeJWT(
		user.ID,
		cfg.secret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making JWT")
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making rToken")
		return
	}

	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 1440),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating rToken record")
	}

	respondWithJSON(w, http.StatusOK, response{
		taggedUser: taggedUser{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		},
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
	})
}
