package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getJWTSecret())

// Claims represents the JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// getJWTSecret retrieves the JWT secret from environment or uses a default for development
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default secret for development only
		// In production, JWT_SECRET environment variable MUST be set
		secret = "youtube-clone-secret-key-change-in-production"
		// Log a warning in production environments
		if os.Getenv("GO_ENV") == "production" {
			panic("JWT_SECRET environment variable must be set in production")
		}
	}
	return secret
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID int, username, role string) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken generates a new token from an existing valid token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same claims
	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}
