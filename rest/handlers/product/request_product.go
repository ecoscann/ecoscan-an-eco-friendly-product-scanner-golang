package product

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
)

func (h *ProductHandler) ReqProduct(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) //10mb max allowed
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	barcode := r.FormValue("barcode")
	name := r.FormValue("name")
	brandName := r.FormValue("brandname")
	userID := 123 //for jwt verification

	if barcode == "" || name == "" {
		http.Error(w, "Barcode and Name is required", http.StatusBadRequest)
		return
	}

	file, imgMetaData, err := r.FormFile("productImage")
	if err != nil {
		http.Error(w, "Image file misisng", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageURL, err := uploadToCloud(file, imgMetaData)
	if err != nil {
		log.Println("Error uploading image: ", err)
		http.Error(w, "Image file misisng", http.StatusBadRequest)
		return
	}

	tx, err := h.DB.Beginx()
	if err != nil {
		log.Println("Error starting db transaction: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	reqQuery := `INSERT INTO product_requests (barcode, name, brand_name, user_id, image_url)
        VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.Exec(reqQuery, barcode, name, brandName, userID, imageURL)
	if err != nil {
		log.Println("Product insertion request failed: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//user submit req : +10 point added to his profile
	pointsQuery := `UPDATE users SET points = points + 10 WHERE id = $1`
	_, err = tx.Exec(pointsQuery, userID)
	if err != nil {
		log.Println("Error: failed to update user points: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error: failed to commit transaction: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{"message": "Request submitted successfully"})

}

func uploadToCloud(file multipart.File, imgMetaData *multipart.FileHeader) (string, error) {
	return "", nil
}
