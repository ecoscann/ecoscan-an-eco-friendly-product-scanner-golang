package product

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"

    "ecoscan.com/repo"
)

func (h *ProductHandler) generateMotivationalMessage(product repo.Product, score int) string {
    apiKey := os.Getenv("OPENROUTER_API_KEY") // set this in Render
    if apiKey == "" {
        return "Choosing eco-friendly products helps reduce waste and protect the planet!"
    }

    // Chat-style input
    messages := []map[string]string{
        {
            "role": "user",
            "content": fmt.Sprintf(
                "Write a short, uplifting eco-friendly motivational message (max 2 sentences) for a user buying %s by %s. Eco Score: %d. Make it sound inspiring and personal. Encourage them to keep using ecoScanAi.",
                product.Name, product.BrandName, score,
            ),
        },
    }

    payload := map[string]interface{}{
        "model":    "mistralai/voxtral-small-24b-2507",
        "messages": messages,
    }

    bodyBytes, _ := json.Marshal(payload)

    url := "https://openrouter.ai/api/v1/chat/completions"

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("HTTP-Referer", "https://yourapp.com") // optional
    req.Header.Set("X-Title", "ecoScanAi")                // optional

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("OpenRouter API error: %v", err)
        return "Your choice makes a positive impact on the environment!"
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)
    log.Printf("OpenRouter status: %d", resp.StatusCode)
    log.Println("OpenRouter raw response:", string(respBody))

    // Parse response
    var result struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }
    if err := json.Unmarshal(respBody, &result); err == nil {
        if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
            return result.Choices[0].Message.Content
        }
    }

    return "Thanks for choosing sustainable products — together we reduce waste!"
}
