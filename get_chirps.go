package main

import (
	"net/http"
)

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirps")
		return
	}

	// allocate a slice of tagged chirps and map the queried chirps onto the elements
	returnChirps := make([]tagChirp, len(chirps))
	for i, chirp := range chirps {
		returnChirps[i] = tagChirp{
			ID:         chirp.ID,
			Created_At: chirp.CreatedAt,
			Updated_At: chirp.UpdatedAt,
			Body:       chirp.Body,
			User_ID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, returnChirps)
}
