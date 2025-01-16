package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/per1Peteia/chirpy/internal/database"
)

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	// collect sort order string
	sortOrder := r.URL.Query().Get("sort")

	// allow for optional query parameter /api/chirps?author_id=
	if id := r.URL.Query().Get("author_id"); id != "" {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error parsing query id")
			return
		}
		chirps, err := cfg.dbQueries.GetAllChirpsByID(r.Context(), parsedID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error getting chirps by id")
			return
		}
		returnChirps := returnPopulatedChirps(chirps, sortOrder)
		respondWithJSON(w, http.StatusOK, returnChirps)
		return
	}
	// otherwise return all chirps for all authors
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirps")
		return
	}
	returnChirps := returnPopulatedChirps(chirps, sortOrder)

	respondWithJSON(w, http.StatusOK, returnChirps)
}

func returnPopulatedChirps(chirps []database.Chirp, order string) []tagChirp {
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
	// depending on order argument select for sort order
	switch order {
	case "asc":
		sort.Slice(returnChirps, func(i, j int) bool {
			return returnChirps[i].Created_At.Before(returnChirps[j].Created_At)
		})
	case "desc":
		sort.Slice(returnChirps, func(i, j int) bool {
			return returnChirps[i].Created_At.After(returnChirps[j].Created_At)
		})
	default:
		sort.Slice(returnChirps, func(i, j int) bool {
			return returnChirps[i].Created_At.Before(returnChirps[j].Created_At)
		})
	}

	return returnChirps
}
