package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"database/sql"
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
	IsChirpyRed 	bool		`json:"is_chirpy_red"`
}

//
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
		sendJsonResponse(w, 400, map[string]string{"error": "Invalid user request payload"})
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
		IsChirpyRed: user.IsChirpyRed,
	})
}

//
func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {
	platform, found := os.LookupEnv("PLATFORM")
	if !found || platform != "dev" {
		sendJsonResponse(w, http.StatusForbidden, map[string]string{"error": "Access restricted in non-development enviroments"})
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error resetting database: %s", err)
		sendJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete all users"})
		return
	}
	sendJsonResponse(w, http.StatusOK, map[string]string{"status": "All users deleted successfully"})
}

//
func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 		string `json:"email"`
		Password 	string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		sendJsonResponse(w, 400, map[string]string{"error": "Invalid user request payload"})
		return 
	}

	if params.Email == "" || params.Password == "" {
		sendJsonResponse(w, 400, map[string]string{"error": "Required fields are empty"})
	}

	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error: %v", err)
		sendJsonResponse(w, 401, map[string]string{"error": "Failed to get users token"})
		return
	}

	userUUID, err := auth.ValidateJWT(userToken, cfg.secret)
	if err != nil {
		log.Printf("Error user not authorized: %v", err)
		sendJsonResponse(w, 401, map[string]string{"error": "Unauthorized"})
		return
	}

	newPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing users password: %v", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to create users password"})
		return
	}

	updatedUser, err := cfg.db.UpdateUsersCredentials(r.Context(), database.UpdateUsersCredentialsParams{
		Email: params.Email,
		HashedPassword: newPassword,
		ID: userUUID,
	})
	if err != nil {
		log.Printf("Error adding users new credentials: %v", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to update credentials"})
		return 
	}

	sendJsonResponse(w, 200, User{
		ID: updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email: updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})
}

//
func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		sendJsonResponse(w, 400, map[string]string{"error": "Invalid user request payload"})
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error: %v", err)
		sendJsonResponse(w, 401, map[string]string{"error": "Failed to get API key"})
		return
	}
	if apiKey != cfg.polkaKey {
		log.Print("Invalid API Key")
		sendJsonResponse(w, 401, map[string]string{"error": "Invalid API key"})
		return
	}

	if params.Event != "user.upgraded" {
		log.Printf("Ignored event: %s is not 'user.upgraded'", params.Event)
		sendJsonResponse(w, 204, nil)
		return
	}

	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), database.UpgradeUserToChirpyRedParams{
		ID: params.Data.UserID,
		IsChirpyRed: true,
	})
	if err == sql.ErrNoRows {
		log.Printf("User not found %v", params.Data.UserID)
		sendJsonResponse(w, 404, map[string]string{"error": "User not found"})
		return
	}

	if err != nil {
		log.Printf("Database error while upgrading user: %v",  err)
		sendJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	log.Print("User was upgraded successfully")
	sendJsonResponse(w, 204, nil)
	return 
}
