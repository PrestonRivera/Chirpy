package main

import (
	"net/http"
	"log"
)

func (cfg *apiConfig)handlerReset(w http.ResponseWriter, req *http.Request) {
	log.Printf("Reset handler called. Method: %s, Path: %s", req.Method, req.URL.Path)
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0")) 
}