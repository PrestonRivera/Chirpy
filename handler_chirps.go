package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID 			uuid.UUID 	`json:"id"`
	Created_at 	time.Time 	`json:"created_at"`
	Updated_at 	time.Time 	`json:"updated_at"`
	Body 		string 		`json:"body"`
	User_id 	uuid.UUID 	`json:"user_id"`
}

//
func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to handle user request"})
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting authorization header: %v", err)
		sendJsonResponse(w, 401, map[string]string{"error": "Failed to get user token"})
		return 
	}

	userUUID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("User not authorized: %v", err)
		sendJsonResponse(w, 401, map[string]string{"error": "Unauthorized"})
		return 
	}

	chirp := isChirpValid(params.Body)
	if len(chirp) < 1 {
		sendJsonResponse(w, 400, nil)
		return
	}
	newChirp, err := cfg.db.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body: chirp,
		UserID: userUUID,
	})
	if err != nil {
		log.Printf("Error creating new chirp: %s", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to create new chirp"})
		return
	}
	sendJsonResponse(w, 201, Chirp{
		ID: newChirp.ID,
		Created_at: newChirp.CreatedAt,
		Updated_at: newChirp.UpdatedAt,
		Body: newChirp.Body,
		User_id: newChirp.UserID,
	})
}

//
func isChirpValid(chirp string) string {
	if len(chirp) > 140 {
		return ""
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	words := strings.Split(chirp, " ")
	
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, exists := badWords[loweredWord]; exists {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

//
func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Failed to get Chirps from database: %s", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to retrieve list of chirps"})
		return 
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID: 		dbChirp.ID,
			Created_at: dbChirp.CreatedAt,
			Updated_at: dbChirp.UpdatedAt,
			Body: 		dbChirp.Body,
			User_id: 	dbChirp.UserID,
		}
	}	
	sendJsonResponse(w, 200, chirps)
}

//
func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {
	reqID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(reqID)
	if err != nil {
		log.Printf("Chirp UUID is invalid: %s", err)
		sendJsonResponse(w, 400, map[string]string{"error": "Invalid chirp ID"})
		return
	}

	chirp, err := cfg.db.GetSingleChirp(r.Context(), parsedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Chirp not found for chirpID: %v", parsedID)
			sendJsonResponse(w, 404, map[string]string{"error": "Chirp not found"})
			return 
		}
		log.Printf("Database error for chirpID: %v, %s", parsedID, err)
		sendJsonResponse(w, 500, map[string]string{"error": "Internal server error"})
		return
	}
	sendJsonResponse(w, 200, Chirp{
		ID: chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body: chirp.Body,
		User_id: chirp.UserID,
	})
}

//
func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	reqID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(reqID)
	if err != nil {
		log.Printf("Chirp UUID is invalid")
		sendJsonResponse(w, 400, map[string]string{"error": "Invalid chirp ID"})
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting authorization header: %v", err)
		sendJsonResponse(w, 401, map[string]string{"error": "Failed to get user token"})
		return
	}

	userUUID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("User not authorized")
		sendJsonResponse(w, 403, map[string]string{"error": "User is not authorized"})
		return
	}

	chirp, err := cfg.db.GetSingleChirp(r.Context(), parsedID)
	if err != nil {
		log.Printf("Chirp not found for chirpID: %v", parsedID)
		sendJsonResponse(w, 404, map[string]string{"error": "Chirp not found"})
		return
	}

	if userUUID != chirp.UserID {
		sendJsonResponse(w, 403, map[string]string{"error": "User not authorized to delete chirp"})
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), parsedID)
	if err != nil {
		log.Printf("Error deleting users chirp: %v", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Internal server error"})
		return
	}
	sendJsonResponse(w, 204, nil)
}