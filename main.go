package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	// An HTTP request multiplexer
	// It matches the URL of each incoming request against a
	// list of registered patterns and calls the handler for
	// the pattern that most closely matches the URL.
	mux := http.NewServeMux()
	// When mux hasn't been assigned any handlers it simply returns 404

	// 1. Create a file server handler.
	// `http.Dir(".")` creates a file system rooted at the "." path.
	fileServer := http.FileServer(http.Dir("."))

	// 2. Register the file server to handle requests for the root path.
	// http.Handle("/", ...) tells the multiplexer to use the fileServer
	// for any path that starts with `/`.
	mux.Handle("/", fileServer)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

// Command to build and run the server:
// go build -o out && ./out
