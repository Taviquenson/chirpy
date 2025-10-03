package main

import (
	"net/http"

	"github.com/Taviquenson/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	// Authenticate refresh token
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find refresh token", err)
		return
	}

	dbUser, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user from refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, auth.TokenAcessExpirationTime)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't create JSON Web Token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	// Authenticate refresh token
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find refresh token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
