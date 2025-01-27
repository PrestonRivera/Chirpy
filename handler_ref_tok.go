package main

import (
	"Chirpy/internal/auth"
	"net/http"
	"time"
)

type RefreshResponse struct {
	Token string `json:"token"`
}

//
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	reqToken,  err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendJsonResponse(w, 401, map[string]string{"error": "Invalid token"})
		return 
	}

	refToken, err := cfg.db.GetRefreshToken(r.Context(), reqToken)
	if err != nil {
		sendJsonResponse(w, 401, map[string]string{"error": "Invalid token"})
		return
	}

	if time.Now().After(refToken.ExpiresAt) {
		sendJsonResponse(w, 401, map[string]string{"error": "Invalid token"})
		return 
	}

	if refToken.RevokedAt.Valid {
		sendJsonResponse(w, 401, map[string]string{"error": "Invalid token"})
		return 
	}

	jwtToken, err := auth.MakeJWT(refToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to access token"})
		return 
	}
	sendJsonResponse(w, 200, RefreshResponse{
		Token: jwtToken,
	})
}

//
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendJsonResponse(w, 401, map[string]string{"error": "Invalid token"})
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), reqToken)
	if err != nil {
		sendJsonResponse(w, 500, map[string]string{"error": "Failed to revoke token"})
		return
	}
	sendJsonResponse(w, 204, nil)
}