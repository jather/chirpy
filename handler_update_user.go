package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jather/chirpy/internal/auth"
	"github.com/jather/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := auth.ValidateJWT(token, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	param := parameters{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect input")
		return
	}
	hashedPassword, err := auth.HashPassword(param.Password)
	if err != nil {
		log.Printf("error while hashing password, %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	updatedUser, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{ID: user, Email: param.Email, HashedPassword: hashedPassword})
	if err != nil {
		log.Printf("database error, %s", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resp := User{Id: updatedUser.ID, CreatedAt: updatedUser.CreatedAt, UpdatedAt: updatedUser.UpdatedAt, Email: updatedUser.Email, IsChirpyRed: updatedUser.IsChirpyRed.Bool}
	respondWithJson(w, http.StatusOK, resp)

}
