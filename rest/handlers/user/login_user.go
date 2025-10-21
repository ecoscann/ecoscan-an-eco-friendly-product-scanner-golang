
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


type LoginResponse struct {
	Message      string    `json:"message"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         repo.User `json:"user"`
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user repo.User
	query := `SELECT * FROM users WHERE email = $1`
	err := h.DB.Get(&user, query, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			log.Printf("Database error finding user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	expiryTime := time.Now().Add(time.Hour * 24 * 7)
	_, err = h.DB.Exec(`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token = $2, expires_at = $3`, user.ID, refreshToken, expiryTime)
	if err != nil {
		log.Printf("Failed to save refresh token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = "" 

	loginResponse := LoginResponse{
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user, 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}