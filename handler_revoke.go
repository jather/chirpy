package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/jather/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := cfg.db.GetUserFromRefreshToken(req.Context(), token)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		} else {
			log.Printf("database error, %s", err)
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	err = cfg.db.RevokeToken(req.Context(), userID)
	if err != nil {
		log.Printf("database error, %s", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
