package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	dbchirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		fmt.Printf("database error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
	resp := make([]Chirp, len(dbchirps))
	for i, chirp := range dbchirps {
		resp[i] = Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID}
	}
	respondWithJson(w, http.StatusOK, resp)
}
