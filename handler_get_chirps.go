package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	authorIDString := req.URL.Query().Get("author_id")
	if authorIDString == "" {
		dbchirps, err := cfg.db.GetChirps(req.Context())
		if err != nil {
			log.Printf("database error: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		resp := make([]Chirp, len(dbchirps))
		for i, chirp := range dbchirps {
			resp[i] = Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID}
		}
		respondWithJson(w, http.StatusOK, resp)
	} else {
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid id")
		}
		_, err = cfg.db.GetUser(req.Context(), authorID)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusBadRequest, "user doesn't exist")
				return
			} else {
				log.Printf("database error, %s", err)
				respondWithError(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}
		dbchirps, err := cfg.db.GetChirpsForUser(req.Context(), authorID)
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithJson(w, http.StatusOK, []Chirp{})
				return
			} else {
				log.Printf("database error, %s", err)
				respondWithError(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}
		resp := make([]Chirp, len(dbchirps))
		for i, chirp := range dbchirps {
			resp[i] = Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID}
		}
		respondWithJson(w, http.StatusOK, resp)

	}
}
