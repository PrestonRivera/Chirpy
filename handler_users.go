package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID 				uuid.UUID 	`json:"id"`
	CreatedAt 		time.Time 	`json:"created_at"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
	Email 			string 		`json:"email"`
}


func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
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
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing user password: %s", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to create users password"})
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating new user: %s", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to create new user"})
		return
	}

	sendJsonResponse(w, 201, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}


func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {
	platform, found := os.LookupEnv("PLATFORM")
	if !found || platform != "dev" {
		sendJsonResponse(w, http.StatusForbidden, map[string]string{"error": "Not authorized"})
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error resetting database: %s", err)
		sendJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to restart database"})
		return
	}
	sendJsonResponse(w, http.StatusOK, map[string]string{"status": "All users deleted successfully"})
}