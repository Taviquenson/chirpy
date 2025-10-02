package main

import (
	"encoding/json"
	"errors"
	// "fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Taviquenson/chirpy/internal/auth"
	"github.com/Taviquenson/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Authenticate user
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get bearer token", err)
		return
	}
	userIDFromToken, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't validate JSON Web Token", err)
		return
	}
	// fmt.Println(userIDFromToken)
	// fmt.Println(params)
	// fmt.Println(params.UserID)
	if userIDFromToken != params.UserID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized to post Chirp", err)
		return
	}

	msgCleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   msgCleaned,
		UserID: params.UserID,
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	// making the keys map to struct{} makes it so they take less space/memory
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	msgCleaned := getCleanedBody(body, badWords)
	return msgCleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, exists := badWords[loweredWord]; exists {
			words[i] = "****"
		}
	}
	msgCleaned := strings.Join(words, " ")
	return msgCleaned
}
