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
        "Write a short, uplifting, and creative eco-friendly message (max 2 sentences) for a user buying %s by %s. Eco Score: %d. Make it sound positive, inspiring, and personal. Encourage them to keep using ecoScanAi.",
        product.Name, product.BrandName, score,
    )

    // ✅ New Inference Providers endpoint
    url := "https://router.huggingface.co/hf-inference/models/mistralai/Mistral-7B-Instruct-v0.2"

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
    log.Println("HF raw response:", string(respBody))

    // Try array format
    var arr []map[string]interface{}
    if err := json.Unmarshal(respBody, &arr); err == nil && len(arr) > 0 {
        if text, ok := arr[0]["generated_text"].(string); ok && text != "" {
            return text
        }
    }

    // Try object format
    var obj map[string]interface{}
    if err := json.Unmarshal(respBody, &obj); err == nil {
        if text, ok := obj["generated_text"].(string); ok && text != "" {
            return text
        }
    }

    return "Thanks for choosing sustainable products — together we reduce waste!"
}