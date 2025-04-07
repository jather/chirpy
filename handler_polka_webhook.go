package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jather/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	param := parameters{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect input")
		return
	}
	if param.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = cfg.db.UpdateUserToRed(req.Context(), param.Data.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "not found")
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
