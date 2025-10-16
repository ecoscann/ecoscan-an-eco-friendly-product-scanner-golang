package product

import (
	"net/http"

	"ecoscan.com/repo"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/products/barcode/{barcode}", http.HandlerFunc(GetProduct))
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	var db *sqlx.DB
	var newProduct []repo.Product
	barcode := r.PathValue("barcode")

	query := db.Get(&)



}
