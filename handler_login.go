package main

import (
	"encoding/json"
	"net/http"

	"github.com/jather/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	param := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, 400, "incorrect input")
		return
	}
	user, err := cfg.db.GetUserFromEmail(req.Context(), param.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "email or password incorrect")
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, param.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "email or password incorrect")
		return
	}
	resp := User{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	respondWithJson(w, http.StatusOK, resp)
}
