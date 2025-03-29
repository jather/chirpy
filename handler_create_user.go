package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jather/chirpy/internal/auth"
	"github.com/jather/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	password, err := auth.HashPassword(param.Password)
	if err != nil {
		log.Printf("error while hashing password, %s", err)
		if err == bcrypt.ErrPasswordTooLong {
			respondWithError(w, 400, "password too long")
		} else {
			respondWithError(w, 400, "password invalid")
		}
	}
	dbUser, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{Email: email, HashedPassword: password})
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
	respondWithJson(w, 201, user)
}
