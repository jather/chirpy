package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("chirp_ID"))
	if err != nil {
		log.Printf("invalid ID")
		respondWithError(w, http.StatusBadRequest, "invalid ID")
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), id)
	if err != nil {
		if err != sql.ErrNoRows {
			respondWithError(w, 500, "Internal server error")
		} else {
			respondWithError(w, 404, "Does not exist")
		}
		return
	}
	resp := Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.CreatedAt, Body: chirp.Body, UserID: chirp.UserID}
	respondWithJson(w, http.StatusOK, resp)
}
