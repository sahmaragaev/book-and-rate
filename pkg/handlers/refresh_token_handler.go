package handlers

import (
	"book-and-rate/pkg/auth"
	"book-and-rate/pkg/config"
	"encoding/json"
	"net/http"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var tokenRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&tokenRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Load the configuration
	cfg := config.LoadConfig("config/config.json")

	// Validate the refresh token
	token, err := auth.ValidateToken(tokenRequest.RefreshToken, *cfg)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*auth.Claims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
		return
	}

	// Ensure that claims.UserID is correctly named based on your Claims structure
	newAccessToken, err := auth.GenerateToken(claims.UserId, *cfg)
	if err != nil {
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"accessToken": newAccessToken,
	}
	json.NewEncoder(w).Encode(response)
}
