package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Taviquenson/chirpy/internal/auth"
)

// Placeholder function signature (probably to be relocated)
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password   string `json:"password"`
		Email      string `json:"email"`
		Expiration int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
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
	if !match || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Produce JSON Web Token
	expirationTime := time.Hour
	if params.Expiration > 3600 || params.Expiration == 0 { // params.Expiration will be 0 if unspecified
		expirationTime = time.Hour
	} else {
		expirationTime = time.Duration(params.Expiration) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't produce JSON Web Token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
