package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jather/chirpy/internal/auth"
	"github.com/jather/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	//decode request body
	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters, %v", err)
		respondWithError(w, 500, "Internal server error")
		return
	}
	//check if has valid JWT
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error while getting bearer token, %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtsecret)
	if err != nil {
		log.Printf("error while validatind jwt, %s", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	//Check length of chirp
	if len(params.Body) > 140 {
		log.Print("chirp too long")
		respondWithError(w, http.StatusBadRequest, "Chirp cannot be more than 140 characters")
	} else {
		params.Body = cleanText(params.Body)
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{Body: params.Body, UserID: userID})
	if err != nil {
		log.Printf("database error, %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
	resp := Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID}
	respondWithJson(w, http.StatusCreated, resp)
}

func cleanText(s string) string {
	badWords := map[string]struct{}{"kerfuffle": {}, "fornax": {}, "sharbert": {}}
	words := strings.Fields(s)
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
