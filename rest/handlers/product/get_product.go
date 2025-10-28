// rest/handlers/product/get.go (সংশোধিত GetProduct ফাংশন)
package product

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ecoscan.com/logic" 
	"ecoscan.com/repo"
)

type ProductResponse struct {
	Product      repo.Product   `json:"product"`
	Score        int            `json:"score"`
	ScoreRating  string         `json:"score_rating"`
	Alternatives []repo.Product `json:"alternatives"`
}


func getScoreRating(score int) string {
    if score == 0 { 
        return "Not Rated" 
    }
	if score <= 30 { return "High Impact" }
	if score <= 60 { return "Moderate Impact" }
	if score <= 80 { return "Good Choice" }
	return "Excellent Choice"
}


func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json") // Set header early

	var mainProduct repo.Product 
	barcode := r.PathValue("barcode")

   
	queryMain := `SELECT * FROM products WHERE barcode = $1;`
	err := h.DB.Get(&mainProduct, queryMain, barcode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"message": "Product not found"}`, http.StatusNotFound)
		} else {
			log.Printf("Database error fetching product: %v", err)
			http.Error(w, `{"message": "Internal server error reading product"}`, http.StatusInternalServerError)
		}
		return
	}

   
	productScore := int(logic.CalculateScore(mainProduct))
	scoreRating := getScoreRating(productScore)
    log.Printf("Calculated score for barcode %s: %d (%s)", barcode, productScore, scoreRating) 

	
	var alternatives []repo.Product
	queryAlt := `
		SELECT * FROM products
		WHERE category = $1 AND id != $2 AND (price < $3 OR packaging_material IN ('glass', 'paper', 'none', 'compostable_paper', 'cardboard'))
		ORDER BY price ASC, packaging_material ASC
		LIMIT 4
	`
	err = h.DB.Select(&alternatives, queryAlt, mainProduct.Category, mainProduct.ID, mainProduct.Price)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { 
		log.Printf("Could not find alternatives for product ID %d: %v", mainProduct.ID, err)
       
	}

	
	response := ProductResponse{
		Product:      mainProduct,  
		Score:        productScore, 
		ScoreRating:  scoreRating,  
		Alternatives: alternatives, 
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Printf("Error encoding response: %v", err)
        
    }
}