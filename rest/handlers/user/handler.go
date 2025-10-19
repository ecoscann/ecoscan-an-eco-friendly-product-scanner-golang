package user

import "github.com/jmoiron/sqlx"

// database operations and access to connect
type UserHandler struct {
	DB *sqlx.DB
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{
		DB: db,
	}
}
