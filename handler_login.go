package main

import (
	"Chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"
)

//
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 	string `json:"password"`
		Email 		string `json:"email"`
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
	sendJsonResponse(w, 200, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}