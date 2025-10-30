package product

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "ecoscan.com/repo"
)

// generateMotivationalMessage calls Hugging Face Inference API
// to generate a short eco-friendly motivational message.
func (h *ProductHandler) generateMotivationalMessage(product repo.Product, score int) string {
    apiKey := os.Getenv("HF_API_KEY")
    if apiKey == "" {
        return "Choosing eco-friendly products helps reduce waste and protect the planet!"
    }

    prompt := fmt.Sprintf(
        "Write a short, motivating message (1-2 sentences) for a user buying %s by %s. Eco Score: %d. Highlight reduced waste and sustainability. Encourage them to keep using ecoScanAi.",
        product.Name, product.BrandName, score,
    )

    url := "https://api-inference.huggingface.co/models/mistralai/Mistral-7B-Instruct-v0.2"

    body := map[string]string{"inputs": prompt}
    bodyBytes, _ := json.Marshal(body)

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("HF API error: %v", err)
        return "Your choice makes a positive impact on the environment!"
    }
    defer resp.Body.Close()

    respBody, _ := ioutil.ReadAll(resp.Body)

    var result []map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        log.Printf("Error parsing Hugging Face response: %v", err)
        return "Every eco-friendly purchase counts toward a greener future!"
    }

    if len(result) > 0 {
        if text, ok := result[0]["generated_text"].(string); ok {
            return text
        }
    }

    return "Thanks for choosing sustainable products — together we reduce waste!"
}
