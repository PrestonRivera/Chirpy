package main

import (
	"Chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	ID 			uuid.UUID 	`json:"ID"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	Email 		string 		`json:"email"`
	Token 		string 		`json:"token"`
}

//
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 	string  `json:"password"`
		Email 		string  `json:"email"`
		ExpiresIn	int		`json:"expires_in_seconds"`
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

	expiresIn := 3600
	if params.ExpiresIn > 0 {
		if params.ExpiresIn > 3600 {
			expiresIn = 3600
		} else {
			expiresIn = params.ExpiresIn
		}
	}
	duration := time.Duration(expiresIn) * time.Second
	userToken, err := auth.MakeJWT(user.ID, cfg.secret, duration)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to create token"})
		return
	}

	sendJsonResponse(w, 200, LoginResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: userToken,
	})
}