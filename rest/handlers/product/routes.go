package product

import (
	"encoding/json"
	"net/http"

	"ecoscan.com/repo"
	"github.com/jmoiron/sqlx"
)

type ProductHandler struct{
	DB *sqlx.DB
}

func NewProductHandler(db *sqlx.DB) *ProductHandler{
	return &ProductHandler{
		DB: db,
	}
}

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/products/barcode/{barcode}", http.HandlerFunc(h.GetProduct))
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct repo.Product
	barcode := r.PathValue("barcode")

	query := `SELECT * FROM products WHERE barcode = $1;
 `

	err := h.DB.Get(&newProduct, query, barcode)
	if err != nil {
		http.Error(w, "No Product Found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(newProduct)

}
