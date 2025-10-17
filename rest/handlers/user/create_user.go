package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Points       int       `json:"points" db:"points"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
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

	var newUser User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
	var usr User
	err = h.DB.Get(&usr, query, newUser.Name, newUser.Email, newUser.PasswordHash)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usr)
}
