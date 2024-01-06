package middleware

import (
	"book-and-rate/pkg/auth"
	"book-and-rate/pkg/config"
	"net/http"
	"strings"
)

// AuthenticationMiddleware verifies the JWT token
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		cfg := config.LoadConfig("./config/config.json")
		token, err := auth.ValidateToken(tokenString, *cfg)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
