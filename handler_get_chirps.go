package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	dbchirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		log.Printf("database error %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
	resp := make([]Chirp, len(dbchirps))
	for i, chirp := range dbchirps {
		resp[i] = Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID}
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error while marshalling struct, %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
	w.Write(payload)
}
