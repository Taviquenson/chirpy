package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, req *http.Request) {
	// Check if a user_id/author_id was included in the request to decide whether to return
	// all chirps or just those authored by the provided user
	s := req.URL.Query().Get("author_id")
	// s is a string that contains the value of the author_id query parameter
	// if it exists, or an empty string if it doesn't
	if s != "" {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}

		chirpsUser, err := cfg.db.GetChirpsByUser(req.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't find chirps for user ID", err)
			return
		}
		chirpsForJSON := []Chirp{}
		for _, chirpDB := range chirpsUser {
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
		return
	}

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

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, req *http.Request) {
	chirpIDString := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}
