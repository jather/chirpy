package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("error while decoding json")
		respondWithError(w, http.StatusBadRequest, "Incorrect input")
		return
	}
	email := param.Email
	if email == "" {
		log.Printf("Email field empty")
		respondWithError(w, http.StatusBadRequest, "Email field cannot be empty")
		return
	}
	dbUser, err := cfg.db.CreateUser(req.Context(), email)
	if err != nil {
		log.Printf("error during database operation")
		respondWithError(w, http.StatusBadRequest, "Internal server error")
		return
	}
	user := User{
		Id:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("database error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	respondWithJson(w, 201, data)
}
