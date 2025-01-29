package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	Email 		string 		`json:"email"`
	IsChirpyRed bool		`json:"is_chirpy_red"`
	Token 		string 		`json:"token"`
	RefreshTok  string		`json:"refresh_token"`
}

//
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 	string  `json:"password"`
		Email 		string  `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to handle user request"})
		return
	}

	user, err := cfg.db.FindUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		sendJsonResponse(w, 401, map[string]string{"error": "Incorrect email or password"})
		return 
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		sendJsonResponse(w, 401, map[string]string{"error": "Incorrect email or password"})
		return
	}

	userToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to create token"})
		return
	}

	refTok, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to create refresh token"})
		return 
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour)
	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refTok,
		UserID: user.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		log.Printf("Error storing refresh token")
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to store user refresh token"})
		return 
	}

	sendJsonResponse(w, 200, LoginResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
		Token: userToken,
		RefreshTok: refTok,
	})
}