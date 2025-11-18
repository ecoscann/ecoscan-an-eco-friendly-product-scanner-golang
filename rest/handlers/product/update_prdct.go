package product

import (
	"encoding/json"
	"net/http"

	"ecoscan.com/repo"
)

func (h *ProductHandler) UpdatePrd(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var reqPrd repo.Product
	json.NewDecoder(r.Body).Decode(&reqPrd)

	queryUpdate := `
    UPDATE products
    SET image_url = $1
    WHERE barcode = $2;`

_, err := h.DB.Exec(queryUpdate, reqPrd.ImageURL, reqPrd.Barcode)
if err != nil {
    http.Error(w, "Failed updating img url", http.StatusInternalServerError)
	return
	}	

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reqPrd)

}