package main

import (
	"encoding/json"
	"net/http"

	"github.com/Taviquenson/chirpy/internal/auth"
	"github.com/Taviquenson/chirpy/internal/database"
)

// Placeholder function signature (probably to be relocated)
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Decode request parameters
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Get database user according to request's email
	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Incorrect email or password", err)
		return
	}
	// Validate request password with stored database password_hash
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if (!match) || (err != nil) {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Produce JSON Web Token
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, auth.TokenAcessExpirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JSON Web Token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}
	_, err = cfg.db.CreateRefreshToken(req.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
