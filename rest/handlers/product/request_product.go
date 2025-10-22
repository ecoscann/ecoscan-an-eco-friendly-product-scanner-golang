
package product

import (
	"context" 
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"ecoscan.com/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func (h *ProductHandler) ReqProduct(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) //10mb max allowed
	if err != nil {
		log.Printf("ERROR parsing multipart form: %v", err)
		http.Error(w, "Could not parse request data", http.StatusBadRequest)
		return
	}

	barcode := r.FormValue("barcode")
	name := r.FormValue("name")
	brandName := r.FormValue("brandname") 
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		log.Println("ERROR: Could not get user ID from context")
		http.Error(w, "User authentication error", http.StatusInternalServerError)
		return
	}

	if barcode == "" || name == "" {
		http.Error(w, "Barcode and Name are required", http.StatusBadRequest)
		return
	}

	file, imgMetaData, err := r.FormFile("productImage")
	if err != nil {
		log.Printf("ERROR getting image file: %v", err)
		http.Error(w, "Image file is missing or invalid", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageURL, err := uploadToCloud(file, imgMetaData)
	if err != nil {
		log.Printf("ERROR uploading image: %v", err)
		http.Error(w, "Could not upload image", http.StatusInternalServerError)
		return 
	}

	tx, err := h.DB.Beginx()
	if err != nil {
		log.Printf("ERROR starting db transaction: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	reqQuery := `INSERT INTO product_requests (barcode, name, brand_name, user_id, image_url)
		VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.Exec(reqQuery, barcode, name, brandName, userID, imageURL)
	if err != nil {
		log.Printf("ERROR inserting product request: %v", err)
		http.Error(w, "Failed to save request", http.StatusInternalServerError)
		return
	}

	pointsQuery := `UPDATE users SET points = points + 10 WHERE id = $1`
	_, err = tx.Exec(pointsQuery, userID)
	if err != nil {
		log.Printf("ERROR updating user points: %v", err)
		
	}


	if err := tx.Commit(); err != nil {
		log.Printf("ERROR committing transaction: %v", err)
		http.Error(w, "Failed to finalize request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Request submitted successfully"})
}


func uploadToCloud(file multipart.File, imgMetaData *multipart.FileHeader) (string, error) {
	cnf := config.GetConfig()
	cldUrl := cnf.CloudinaryURL
	if cldUrl == "" {
		log.Println("Cloudinary URL not configured")
		return "", fmt.Errorf("cloudinary URL not configured") 
	}

	cld, err := cloudinary.NewFromURL(cldUrl)
	if err != nil {
		log.Printf("Failed to init Cloudinary: %v", err)
		return "", err
	}

	ctx := context.Background()
	uploadParams := uploader.UploadParams{
		Folder: "ecoscan_products",
	}

	uploadResult, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		return "", err
	}
	return uploadResult.SecureURL, nil
}
