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
	Product     repo.Product
	Alernatives []repo.Product
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {

	var mainProduct repo.Product
	barcode := r.PathValue("barcode")

	query := "SELECT * FROM products WHERE barcode=$1"

	err := h.DB.Get(&mainProduct, query, barcode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			log.Println(w, "Database error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// alternative products
	var alternatives []repo.Product
	queryAlt := `SELECT * FROM products 
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
	if err != nil{ //its the not found error
		log.Println("Couldn't find alternatives for product ID: ", mainProduct.ID)
		alternatives = []repo.Product{}
	}

	// cooking the response

	response := ProductResponse{
		Product: mainProduct,
		Alernatives: alternatives,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	

}
