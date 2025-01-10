package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/per1Peteia/chirpy/internal/database"
)

type apiConfig struct {
	dbQueries      *database.Queries
	fileserverHits atomic.Int32
}

func main() {
	// load .env into my environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// open connection to database
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	// store queries code in config so handlers can access it
	dbQueries := database.New(db)
	apiCfg := &apiConfig{dbQueries: dbQueries}

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fileServer := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/app/", http.StripPrefix("/app", fileServer))

	mux.HandleFunc("GET /api/healthz", readinessHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.srvHitsHandler)

	mux.HandleFunc("POST /admin/reset", apiCfg.resetSrvHitsHandler)

	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)

	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	log.Printf("serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}
