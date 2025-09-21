package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	// An HTTP request multiplexer
	// It matches the URL of each incoming request against a
	// list of registered patterns and calls the handler for
	// the pattern that most closely matches the URL.
	mux := http.NewServeMux()
	// When mux hasn't been assigned any handlers it simply returns 404

	// Create a file server handler.
	fileServer := http.FileServer(http.Dir(filepathRoot))

	// Register the file server to handle requests at /app/
	mux.Handle("/app/", http.StripPrefix("/app/", fileServer))
	// http.FileServer serves files relative to the directory
	// you give it (http.Dir(filepathRoot)) and uses the full
	// URL path to find files. Without stripping the "/app"
	// prefix, it will look for files under filepathRoot/app/...
	// which likely doesn’t exist.

	// Readiness endpoints are used by external systems to check
	// if our server is ready to receive traffic.
	// The endpoint should be accessible at the /healthz path
	// using any HTTP method.
	mux.HandleFunc("/healthz", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// Command to build and run the server:
// go run .

// Command to compile a binary and run it
// go build -o out && ./out
