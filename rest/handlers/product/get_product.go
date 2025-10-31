package product

import (
	"encoding/json"
	"net/http"

	"ecoscan.com/logic"
	"ecoscan.com/repo"
)

type ProductResponse struct {
	Product      repo.Product   `json:"product"`
	Score        int            `json:"score"`
	ScoreRating  string         `json:"score_rating"`
	Alternatives []repo.Product `json:"alternatives"`
	Message      string         `json:"message"`
}

func getScoreRating(score int) string {
	if score <= 0 {
		return "Not Rated"
	}
	if score <= 30 {
		return "High Impact"
	}
	if score <= 60 {
		return "Moderate Impact"
	}
	if score <= 80 {
		return "Good Choice"
	}
	return "Excellent Choice"
}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	barcode := r.PathValue("barcode")
	var mainProduct repo.Product
	query := `SELECT id, barcode, name, brand_name, category, sub_category,
              image_url, price, packaging_material, manufacturing_location, disposal_method
              FROM products WHERE barcode = $1;`
	if err := h.DB.Get(&mainProduct, query, barcode); err != nil {
		http.Error(w, `{"message":"Product not found"}`, http.StatusNotFound)
		return
	}

	score := int(logic.CalculateScore(mainProduct))
	rating := getScoreRating(score)


	resp := ProductResponse{
		Product:      mainProduct,
		Score:        score,
		ScoreRating:  rating,
		Alternatives: []repo.Product{}, 
		Message:      "",
	}
	_ = json.NewEncoder(w).Encode(resp)

	
	go func(p repo.Product, s int) {
		msg := h.generateMotivationalMessage(p, s)
		h.Store.Set(p.Barcode, msg)
	}(mainProduct, score)
}

func (h *ProductHandler) GetProductMessage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    barcode := r.PathValue("barcode")
    if msg, ok := h.Store.Get(barcode); ok {
        _ = json.NewEncoder(w).Encode(map[string]string{"message": msg})
        return
    }

    _ = json.NewEncoder(w).Encode(map[string]string{"message": "Generating eco tip…"})
}
