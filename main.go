package main

import (
	"log"
	"net/http"
)

// run your server:   go build -o out && ./out


func main() {
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))

	s := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %v", err.Error())
	}
}


func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}