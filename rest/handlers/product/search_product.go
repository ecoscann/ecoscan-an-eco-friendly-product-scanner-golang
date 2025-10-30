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
	// conver it to lowercase
	searchQueryLower := strings.ToLower(query)

	var results []repo.Product
	// threshold 0.3 : means 0.3 is better for spelling mistakes
	similarityThreshold := 0.3
	maxResults := 10 

	
	dbQuery := `
		SELECT *
		FROM products
		WHERE similarity(name, $1) > $2 OR lower(name) = $3
		ORDER BY
			(lower(name) = $3) DESC,  -- Exact match (case-insensitive) gets highest priority (true = 1, false = 0)
			similarity(name, $1) DESC -- Then order by similarity score
		LIMIT $4
	`

	err := h.DB.Select(&results, dbQuery, query, similarityThreshold, searchQueryLower, maxResults) // sending query and searchlower both

	if err != nil {
		log.Printf("ERROR searching products by name '%s': %v", query, err)
		http.Error(w, `{"message": "Could not perform search"}`, http.StatusInternalServerError)
		return
	}

	// if exact match found, only exact match will send even there is more results
	if len(results) > 0 && strings.ToLower(results[0].Name) == searchQueryLower {
		log.Printf("Exact match found for '%s', returning only that.", query)
        // If strict exact match is needed and ONLY that one should be returned:
        results = []repo.Product{results[0]}
	} else if len(results) == 0 {
        log.Printf("No products found matching query: '%s'", query)
    } else {
         log.Printf("Returning %d fuzzy matches for query: '%s'", len(results), query)
    }


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results) 
}