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
    apiKey := os.Getenv("OPENROUTER_API_KEY") // set in Render
    if apiKey == "" {
        return "Choosing eco-friendly products helps reduce waste and protect the planet!"
    }

    // Adjust prompt based on eco score
    var prompt string
    if score < 50 {
        prompt = fmt.Sprintf(
            "The user is considering buying %s by %s. Eco Score: %d (low). " +
                "Write a short, supportive eco-friendly message (max 2 sentences). " +
                "Encourage them to try alternative choices with higher eco scores for a better impact. " +
                "Still make it positive and motivating. Mention that buying this saves about %d%% of wastage compared to less eco-friendly options.",
            product.Name, product.BrandName, score, score/2, // simple % calculation
        )
    } else {
        prompt = fmt.Sprintf(
            "The user is buying %s by %s. Eco Score: %d (good). " +
                "Write a short, uplifting eco-friendly motivational message (max 2 sentences). " +
                "Make it inspiring and personal. Mention that buying this saves about %d%% of wastage compared to less eco-friendly options.",
            product.Name, product.BrandName, score, score/2,
        )
    }

    // Chat-style input
    messages := []map[string]string{
        {"role": "user", "content": prompt},
    }

    payload := map[string]interface{}{
        "model":    "openai/gpt-4o",
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
