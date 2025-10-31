package product

import (
	"github.com/jmoiron/sqlx"
)

type ProductHandler struct {
	DB    *sqlx.DB
}

func NewProductHandler(db *sqlx.DB) *ProductHandler {
	return &ProductHandler{
		DB:    db,
	}
}
