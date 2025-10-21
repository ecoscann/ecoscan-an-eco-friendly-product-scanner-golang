package product

import (
	"database/sql" 
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ecoscan.com/repo"
)


type ProductResponse struct {
	Product      repo.Product   `json:"product"`     
	Alternatives []repo.Product `json:"alternatives"` 
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var mainProduct repo.Product
	barcode := r.PathValue("barcode")

	queryMain := `SELECT * FROM products WHERE barcode = $1;`
	err := h.DB.Get(&mainProduct, queryMain, barcode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	
	var alternatives []repo.Product
	queryAlt := `
		SELECT * FROM products 
		WHERE category = $1
		  AND id != $2
		  AND (
			  price < $3 OR
			  packaging_material IN ('glass', 'paper', 'compostable_paper', 'none')
		  )
		ORDER BY price ASC, packaging_material ASC
		LIMIT 4
	`
	err = h.DB.Select(&alternatives, queryAlt, mainProduct.Category, mainProduct.ID, mainProduct.Price)
	if err != nil {
		log.Printf("Could not find alternatives for product ID %d: %v", mainProduct.ID, err)
		alternatives = []repo.Product{} 
	}

	response := ProductResponse{
		Product:      mainProduct, 
		Alternatives: alternatives,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response) 
}