package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("reset is only allowed in dev environement"))
		return
	}

	if err := cfg.dbQueries.ResetUsers(req.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting users")
	}
	cfg.fileserverHits.Store(0)
	w.Write([]byte("hits reset to 0"))
}
