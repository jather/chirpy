package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jather/chirpy/internal/database"
)

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

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	w.Header().Add("Content-Type", "text/json")
	//decode request body
	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters, %v", err)
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"something went wrong"}`))
		return
	}
	//Check length of chirp
	if len(params.Body) > 140 {
		log.Print("chirp too long")
		http.Error(w, "chirp too long", http.StatusBadRequest)
	} else {
		params.Body = cleanText(params.Body)
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{Body: params.Body, UserID: params.UserID})
	if err != nil {
		log.Printf("database error, %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
	resp := Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID}
	//encode response body
	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"something went wrong"}`))
	}
	w.WriteHeader(201)
	w.Write(data)
}
