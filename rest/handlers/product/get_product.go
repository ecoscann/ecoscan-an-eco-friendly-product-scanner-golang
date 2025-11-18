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
    Message string `json:"message"`
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

    var mainProduct repo.Product
    barcode := r.PathValue("barcode")

    queryMain := `
        SELECT id, barcode, name, brand_name, category, sub_category,
               image_url, price, packaging_material, manufacturing_location, disposal_method
        FROM products WHERE barcode = $1;`
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
    log.Printf("Calculated score for main product %s: %d (%s)", barcode, productScore, scoreRating)

    var alternativesData []repo.Product
    queryAlt := `
        WITH sub_alts AS (
    SELECT id, barcode, name, brand_name, category, sub_category,
           image_url, price, packaging_material, manufacturing_location, disposal_method
    FROM products
    WHERE sub_category = $1
      AND id != $2
      AND (price < $3 OR packaging_material IN ('glass','paper','none','compostable_paper','cardboard'))
    ORDER BY price DESC, packaging_material ASC
    LIMIT 4
)
SELECT * FROM sub_alts
UNION ALL
SELECT id, barcode, name, brand_name, category, sub_category,
       image_url, price, packaging_material, manufacturing_location, disposal_method
FROM products
WHERE category = $4   -- same category as main product
  AND id != $2
  AND (price < $3 OR packaging_material IN ('glass','paper','none','compostable_paper','cardboard'))
  AND NOT EXISTS (SELECT 1 FROM sub_alts)   -- only if no sub_category alts found
ORDER BY price DESC, packaging_material ASC
LIMIT 4;

    `
    err = h.DB.Select(&alternativesData, queryAlt, mainProduct.SubCatergory, mainProduct.ID, mainProduct.Price)
    if err != nil && !errors.Is(err, sql.ErrNoRows) {
        log.Printf("Could not find alternatives for product ID %d: %v", mainProduct.ID, err)
    }

    for i := range alternativesData {
        altScore := int(logic.CalculateScore(alternativesData[i]))
        alternativesData[i].Score = altScore
        log.Printf("Calculated score for alternative %s: %d", alternativesData[i].Barcode, altScore)
    }

    
    // if no cache we save into db 
    message := h.generateMotivationalMessage(mainProduct, productScore)


    if barcode == "8941193041031"{
        message= "Bashundhara Paper Towel বেশ ভালো একটি পণ্য। তবে প্লাস্টিকের প্যাকেজিং পরিবেশের জন্য ক্ষতিকর হতে পারে, এটা নিয়ে আমাদের সকলকে সচেতন হতে হবে। আপনি নিচে আমাদের Alternative পণ্যগুলো দেখতে পারেন, যেগুলো পরিবেশবান্ধব এবং এর মাধ্যমে প্রায় ৩৭% মতো বর্জ্য দূষণ কমাতে পারবেন🌱।"
    }
    if barcode == "894110001003"{
        message= "Coca-Cola যেকোনো মুহূর্তকে আর রেফ্রেশিং করে তুলে🌱 এই প্যাকেজিংটা প্লাস্টিক হলেও তুলনামূলকভাবে পরিবেশবান্ধব। আপনি নিচে আমাদের Alternatives পণ্যগুলো দেখতে পারেন। পরিবেশ রক্ষায় এভাব আপনার অবদান রাখুন। 🌱।"
    }
    if barcode == "894110001473"{
        message= "Pepsi প্রতিটি moment-কে করে তোলে আরও lively আর energetic ✨ ক্যান প্যাকেজিং হওয়ায় এটি easily recyclable এবং eco-friendly। আপনার এই conscious choice পরিবেশ রক্ষায় একটি গুরুত্বপূর্ণ step 🌍। আমরা আপনার decision-কে সত্যিই appreciate করি🌱।"
    }

    if barcode == "894110001004"{
        message= "Clemon Lemon Soda প্রতিটি sip-কে করে তোলে আরও refreshing 🍋✨ 250ml Can প্যাকেজিং হওয়ায় এটি super easy to recycle এবং eco-friendly choice। আপনার এই cool decision পরিবেশ রক্ষায় একটি ছোট কিন্তু impactful step 🌍। আমরা আপনার conscious lifestyle-কে সত্যিই appreciate করি🌱।"
    }
    response := ProductResponse{
        Product:      mainProduct,
        Score:        productScore,
        ScoreRating:  scoreRating,
        Alternatives: alternativesData,
        Message: message,
    }

    w.WriteHeader(http.StatusOK)
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Printf("Error encoding response: %v", err)
    }
}