package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, req *http.Request) {
	chirpsDB, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	chirpsForJSON := []Chirp{}
	for _, chirpDB := range chirpsDB {
		// Get the database.Chirp into main.Chirp form (with JSON tags)
		chirp := Chirp{
			ID:        chirpDB.ID,
			CreatedAt: chirpDB.CreatedAt,
			UpdatedAt: chirpDB.UpdatedAt,
			Body:      chirpDB.Body,
			UserID:    chirpDB.UserID,
		}
		chirpsForJSON = append(chirpsForJSON, chirp)
	}

	respondWithJSON(w, http.StatusOK, chirpsForJSON)
}
