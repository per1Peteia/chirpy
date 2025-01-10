package main

import "net/http"

func (cfg *apiConfig) resetSrvHitsHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Write([]byte("hits reset to 0"))
}
