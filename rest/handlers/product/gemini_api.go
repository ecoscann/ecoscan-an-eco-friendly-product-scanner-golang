package product

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"ecoscan.com/repo"
)

func(h *ProductHandler) generateMotivationalMessage(product repo.Product, score int) string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "Choosing eco-friendly products helps reduce waste and protect the planet!"
	}

	prompt := `
    Create a 3 sentences, motivating message for a user who is considering purchasing this product:
    Name: ` + product.Name + `
    Brand: ` + product.BrandName + `
    Eco Score: ` + getScoreRating(score) + `

    The message should highlight how purchasing this product reduces environmental impact
    (e.g., waste reduction, sustainability) and encourage the user to keep using ecoScanAi.
    Keep it under 3 sentences. better response considering everything i told you
    `

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + apiKey

	body := []byte(`{
        "contents": [{
            "parts":[{"text": "` + prompt + `"}]
        }]
    }`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Gemini API error: %v", err)
		return "Your choice makes a positive impact on the environment!"
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("Error parsing Gemini response: %v", err)
		return "Every eco-friendly purchase counts toward a greener future!"
	}

	if candidates, ok := result["candidates"].([]interface{}); ok && len(candidates) > 0 {
		cand := candidates[0].(map[string]interface{})
		if content, ok := cand["content"].(map[string]interface{}); ok {
			if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
				if text, ok := parts[0].(map[string]interface{})["text"].(string); ok {
					return text
				}
			}
		}
	}

	return "Thanks for choosing sustainable products — together we reduce waste!"
}
