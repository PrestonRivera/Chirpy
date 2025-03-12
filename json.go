package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//
func sendJsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	if statusCode == 204 {
		w.WriteHeader(statusCode)
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return 
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

/*
func parseJsonRequest(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	err := json.NewDecoder(r.Body).Decode(destination)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		sendJsonResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return err
	}
	return nil
}*/