package product

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"ecoscan.com/repo"
)

func (h *ProductHandler) SearchProductsByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, `{"message": "Query parameter 'q' is required"}`, http.StatusBadRequest)
		return
	}
	
	searchQueryLower := strings.ToLower(query)

	var results []repo.Product
	similarityThreshold := 0.3
	maxResults := 10

	
	dbQuery := `
		SELECT *
		FROM products
		WHERE similarity(lower(name), $1) > $2 OR lower(name) = $1
		ORDER BY
			(lower(name) = $1) DESC,  
			similarity(lower(name), $1) DESC 
		LIMIT $3
	`

	
	err := h.DB.Select(&results, dbQuery, searchQueryLower, similarityThreshold, maxResults)

	if err != nil {
		
		log.Printf("FATAL ERROR searching products by name '%s': %v", query, err)
		http.Error(w, `{"message": "Could not perform search"}`, http.StatusInternalServerError)
		return
	}

	if len(results) > 0 && strings.ToLower(results[0].Name) == searchQueryLower {
		log.Printf("Exact match found for '%s', returning only that.", query)
		
		results = []repo.Product{results[0]}
	} else if len(results) == 0 {
		log.Printf("No products found matching query: '%s'", query)
	} else {
		log.Printf("Returning %d fuzzy matches for query: '%s'", len(results), query)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}