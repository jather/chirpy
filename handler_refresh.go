package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/jather/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := cfg.db.GetUserFromRefreshToken(req.Context(), token)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no rows?")
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		} else {
			log.Printf("error in database lookup, %s", err)
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	expiresInSeconds := 3600
	jwtToken, err := auth.MakeJWT(userID, cfg.jwtsecret, time.Duration(expiresInSeconds)*time.Second)
	if err != nil {
		log.Printf("error making JWT, %s", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	type refreshResp struct {
		Token string `json:"token"`
	}
	resp := refreshResp{Token: jwtToken}
	respondWithJson(w, http.StatusOK, resp)
}
