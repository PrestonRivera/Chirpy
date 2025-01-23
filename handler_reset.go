package main

import (
	"net/http"
	"log"
)

//
func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	log.Printf("Reset fileserverHits handler called. Method: %s, Path: %s", r.Method, r.URL.Path)
	cfg.fileserverHits.Store(0)
	sendJsonResponse(w, http.StatusOK, map[string]string{"status": "Hits reset to 0"})
}