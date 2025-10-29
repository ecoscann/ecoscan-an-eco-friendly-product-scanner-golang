package product

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"ecoscan.com/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func (h *ProductHandler) ReqProduct(w http.ResponseWriter, r *http.Request) {
	// Parse form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("ERROR parsing multipart form: %v", err)
		http.Error(w, "Could not parse request data", http.StatusBadRequest)
		return
	}

	barcode := r.FormValue("barcode")
	name := r.FormValue("name")
	brandName := r.FormValue("brandname") // front-end should match this key

	// validate required fields
	if strings.TrimSpace(barcode) == "" || strings.TrimSpace(name) == "" {
		http.Error(w, "Barcode and Name are required", http.StatusBadRequest)
		return
	}

	// safe userID extraction from context
	userID, ok := extractUserIDFromContext(r.Context())
	if !ok {
		log.Println("ERROR: Could not get user ID from context")
		http.Error(w, "User authentication error", http.StatusInternalServerError)
		return
	}

	// get file
	file, imgMetaData, err := r.FormFile("productImage")
	if err != nil {
		log.Printf("ERROR getting image file: %v", err)
		http.Error(w, "Image file is missing or invalid", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// upload
	imageURL, err := uploadToCloud(file, imgMetaData)
	if err != nil {
		log.Printf("ERROR uploading image: %v", err)
		http.Error(w, "Could not upload image", http.StatusInternalServerError)
		return
	}

	// DB transaction
	tx, err := h.DB.Beginx()
	if err != nil {
		log.Printf("ERROR starting db transaction: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	// ensure rollback if not committed
	defer func() {
		_ = tx.Rollback()
	}()

	reqQuery := `INSERT INTO product_requests (barcode, name, brand_name, user_id, image_url)
        VALUES ($1, $2, $3, $4, $5)`

	if _, err := tx.Exec(reqQuery, barcode, name, brandName, userID, imageURL); err != nil {
		log.Printf("ERROR inserting product request: %v", err)
		http.Error(w, "Failed to save request", http.StatusInternalServerError)
		return
	}

	pointsQuery := `UPDATE users SET points = points + 10 WHERE id = $1`
	if _, err := tx.Exec(pointsQuery, userID); err != nil {
		log.Printf("ERROR updating user points: %v", err)
		http.Error(w, "Failed to update user points", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("ERROR committing transaction: %v", err)
		http.Error(w, "Failed to finalize request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Request submitted successfully"}); err != nil {
		log.Printf("ERROR writing response: %v", err)
	}
}

// helper:  userID extraction
func extractUserIDFromContext(ctx context.Context) (int64, bool) {
	val := ctx.Value("userID")
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	case string:
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			return parsed, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func uploadToCloud(file multipart.File, imgMetaData *multipart.FileHeader) (string, error) {
	// reset reader to start if possible
	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

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
	publicID := strings.TrimSuffix(imgMetaData.Filename, filepath.Ext(imgMetaData.Filename))
	uploadParams := uploader.UploadParams{
		Folder:   "ecoscan_products",
		PublicID: publicID,
	}

	uploadResult, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		return "", err
	}
	return uploadResult.SecureURL, nil
}
