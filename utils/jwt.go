package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"ecoscan.com/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID int64) (string, error) {
	cnf := config.GetConfig()
	claims := jwt.MapClaims{ //this is payload
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), //exp in 15 min
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // it returns header+payload

	return token.SignedString([]byte(cnf.JWTSecretKey)) //signed is adding the signature last
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)  // a slice that wil store rand string of 32byte
	_, err := rand.Read(b) // rand string loading

	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
