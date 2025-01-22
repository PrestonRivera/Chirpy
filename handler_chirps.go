package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

//
func handlerChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sendJsonResponse(w, 500, map[string]string{"error": "Something went wrong"})
		return
	}

	chirp := isChirpValid(params.Body)
	if len(chirp) < 1 {
		sendJsonResponse(w, 400, nil)
		return
	}
	sendJsonResponse(w, 200, map[string]string{"cleaned_body": chirp})
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