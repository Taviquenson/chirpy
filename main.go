package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Taviquenson/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32 // allows us to safely increment and
	// read an integer value across multiple
	// goroutines (HTTP requests)
	db        *database.Queries
	platform  string
	jwtSecret string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Using a connection URL because it can point to a local or remote database
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	// os.Getenv returns "" if the variable isn’t set.
	// If you build apiCfg before loading env (or the var doesn’t exist),
	// cfg.platform will be "" and your check will treat it as non-dev,
	// blocking the reset.
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("jwtSecret must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL) // connect to db
	if err != nil {
		log.Fatalf("Error opening databse: %s", err)
	}
	dbQueries := database.New(dbConn) // type *database.Queries

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	// An HTTP request multiplexer
	// It matches the URL of each incoming request against a
	// list of registered patterns and calls the handler for
	// the pattern that most closely matches the URL.
	mux := http.NewServeMux()
	// When mux hasn't been assigned any handlers it simply returns 404

	// Create a file server handler.
	fileServerHandler := http.FileServer(http.Dir(filepathRoot))

	// Register the file server to handle requests at /app/
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServerHandler))
	mux.Handle("/app/", fsHandler)
	// http.FileServer serves files relative to the directory
	// you give it (http.Dir(filepathRoot)) and uses the full
	// URL path to find files. Without stripping the "/app"
	// prefix, it will look for files under filepathRoot/app/...
	// which likely doesn’t exist.

	// Readiness endpoints are used by external systems to check
	// if our server is ready to receive traffic.
	// The endpoint should be accessible at the /healthz path
	// using any HTTP method.
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

// Command to build and run the server:
// go run .

// Command to compile a binary and run it
// go build -o out && ./out

// Access server at:
// localhost:8080
