package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"ecoscan.com/repo"
	"ecoscan.com/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// searching db with mail to match

	var user repo.User
	query := `SELECT * FROM users WHERE email = $1`
	err := h.DB.Get(&user, query, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		log.Println("Database  error finding user: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// comparing user old hash pass with new one

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Println("Invalid password attempt for email: ", req.Email)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// generating access tokens

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// generating refresh token

	refreshToken, err := utils.GenerateRefreshToken()

	//  saving refresh token into db

	expiryTime := time.Now().Add(time.Hour * 24*7) // 7 days
	query = `INSERT INTO refresh_tokens (user_id, token, expires_at)VALUES ($1, $2, $3)`
	_, err = h.DB.Exec(query, user.ID, refreshToken, expiryTime)

	if err != nil {
        log.Printf("Failed to save refresh token: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//    sending token to the client

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})

}
