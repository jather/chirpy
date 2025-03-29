package main

import (
	"encoding/json"
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
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		return
	}
	email := param.Email
	if email == "" {
		log.Printf("Email field empty")
		http.Error(w, "Email field is empty", http.StatusBadRequest)
		return
	}
	dbUser, err := cfg.db.CreateUser(req.Context(), email)
	if err != nil {
		log.Printf("error during database operation")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
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
		log.Printf("error during database operation")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	w.Write(data)
}
