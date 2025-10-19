package user

import (
	"encoding/json"
	"net/http"

	"ecoscan.com/repo"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// database operations and access to connect
type UserHandler struct {
	DB *sqlx.DB
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{
		DB: db,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Name, email and password required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		http.Error(w, "Password length must be 8 minimum", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO users (
			name, 
			email, 
			password_hash, 
			created_at, 
			updated_at
		)
		VALUES (
			$1, 
			$2, 
			$3, 
			NOW(), 
			NOW()
		)
		RETURNING id, name, email, password_hash, points, created_at, updated_at
	`
	var newUser repo.User
	err = h.DB.Get(&newUser, query, req.Name, req.Email, req.Password)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
