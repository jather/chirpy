package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jather/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	//check token
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	//check that chirp exists
	chirpID, err := uuid.Parse(req.PathValue("chirp_ID"))
	if err != nil {
		log.Printf("invalid ID")
		respondWithError(w, http.StatusBadRequest, "invalid ID")
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "not found")
			return
		} else {
			log.Printf("database error, %s", err)
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	//check if identity of token matches the chirp's owner
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "only owner of chirp can delete it")
		return
	}
	//delete chirp
	err = cfg.db.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		log.Printf("database error, %s", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
