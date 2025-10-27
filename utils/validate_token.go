package utils

import (
	"fmt"
	"log"

	"ecoscan.com/config"
	"github.com/golang-jwt/jwt/v5"
)

// ValidateAccessToken parses and validates a JWT access token and returns claims.
func ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	cnf := config.GetConfig() 

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key from the loaded config
		return []byte(cnf.JWTSecretKey), nil 
	})

	if err != nil {
		// Log the specific parsing error for debugging
		log.Printf("Error parsing token: %v", err)
		return nil, err // Return the specific error (e.g., "token is expired")
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	// If token.Valid is false for some other reason
	return nil, fmt.Errorf("invalid token")
}
