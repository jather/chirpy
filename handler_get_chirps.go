package main

import (
	"database/sql"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/jather/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	authorIDString := req.URL.Query().Get("author_id")
	var chirps []database.Chirp
	//if author_id query is empty, get all chirps
	if authorIDString == "" {
		allChirps, err := cfg.db.GetChirps(req.Context())
		if err != nil {
			log.Printf("database error: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		chirps = allChirps
	} else {
		// parse author_id query
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid id")
		}
		//check that user exists
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
		//get all chirps for user
		filteredChirps, err := cfg.db.GetChirpsForUser(req.Context(), authorID)
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
		chirps = filteredChirps

	}
	//populate resp payload
	resp := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		resp[i] = Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID}
	}
	sortOrder := req.URL.Query().Get("sort")
	if sortOrder == "desc" {
		sort.Slice(resp, func(i, j int) bool { return resp[i].CreatedAt.After(resp[j].CreatedAt) })
	} else {
		sort.Slice(resp, func(i, j int) bool { return resp[i].CreatedAt.Before(resp[j].CreatedAt) })
	}
	respondWithJson(w, http.StatusOK, resp)
}
