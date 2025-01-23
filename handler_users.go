package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID 			uuid.UUID `json:"id"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	Email 		string `json:"email"`
}


func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sendJsonResponse(w, 500, map[string]string{"error": "something went wrong"})
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	
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