package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	err := cfg.db.ResetUsers(req.Context())
	if err != nil {
		log.Printf("database error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
}
