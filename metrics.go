package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	// Set the Content-Type header to indicate the response body format.
	w.Header().Set("Content-Type", "text/html")
	// Set the HTTP status code.
	w.WriteHeader(http.StatusOK) // Or other status codes like http.StatusNotFound, http.StatusInternalServerError
	// Write the response body.
	w.Write([]byte(fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// Wrap the given handler `next` in a nameless handler so that once
	// middlewareMetricsInc() is called and registers this namesless handler,
	// the count will increase and the `next` handle will execute everytime the
	// corresponding endpoint is reached because it is this nameless handler
	// what will actually run on each request, executing its ServeHTTP method
	// (the part inside {} below), which is  what gets registered to the mux
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
