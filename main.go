package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	log.Printf("serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
