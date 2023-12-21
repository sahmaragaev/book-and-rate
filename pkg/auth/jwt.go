package auth

import (
	"book-and-rate/pkg/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims struct
type Claims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func GenerateToken(userID string, cfg config.Config) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JwtSecret))

	return tokenString, err
}

func GenerateRefreshToken(userID string, cfg config.Config) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JwtSecret))
}

// ValidateToken validates the JWT token
func ValidateToken(tokenString string, cfg config.Config) (*jwt.Token, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtSecret), nil
	})

	return token, err
}
