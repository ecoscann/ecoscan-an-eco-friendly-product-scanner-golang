package product
/* 
import (
	"encoding/json"
	"log"
	"net/http"

	"ecoscan.com/repo"
)

func ReqProduct(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10<<20)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	barcode := r.FormValue("barcode")
	name := r.FormValue("name")
	brandName := r.FormValue("brandname")
	userID := 123 //for jwt verification

	file, imgMetaData, err := r.FormFile("productImage")
	if err != nil{
		http.Error(w, "Image file misisng", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgURL, err := uploadToCloud(file, imgMetaData.Filename)
	if err != nil{
		log.Println("Error uploading image: ", err)
		http.Error(w, "Image file misisng", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Request Created")

}

func uploadToCloud(file multipart.File, imgMetaData *multipart.FileHeader) (string, error){

} */