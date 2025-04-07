package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jather/chirpy/internal/auth"
	"github.com/jather/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type loginResp struct {
		Id           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
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

	expiresInSeconds := 3600

	token, err := auth.MakeJWT(user.ID, cfg.jwtsecret, (time.Duration(expiresInSeconds) * time.Second))
	if err != nil {
		log.Printf("error while creating jwt token, %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error while making refresh token, %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server erro")
	}
	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{Token: refreshToken, UserID: user.ID})
	if err != nil {
		log.Printf("database error, %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
	}
	resp := loginResp{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, IsChirpyRed: user.IsChirpyRed.Bool, Email: user.Email, Token: token, RefreshToken: refreshToken}
	respondWithJson(w, http.StatusOK, resp)
}
