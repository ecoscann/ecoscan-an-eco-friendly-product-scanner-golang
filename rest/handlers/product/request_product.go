package product

import (
	"encoding/json"
	"net/http"
	"time"

	"ecoscan.com/repo"
)

type ProductRequest struct{
	ID        int64     `json:"id" db:"id"`
    Barcode   string    `json:"barcode" db:"barcode"`
    Name      string    `json:"name" db:"name"`
    BrandName string    `json:"brand_name" db:"brand_name"`
    ImageURL  string    `json:"image_url" db:"image_url"`
    UserID    int64     `json:"user_id" db:"user_id"`
    Status    string    `json:"status" db:"status"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func ReqProduct(w http.ResponseWriter, r *http.Request){
	var reqPrd repo.Product
	err := json.NewDecoder(r.Body).Decode(&reqPrd) 
	if err != nil{
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}


}