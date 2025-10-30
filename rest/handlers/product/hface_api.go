package product

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "math/rand"
    "net/http"
    "os"
    "time"

    "ecoscan.com/repo"
)

// Pre-written fallback messages (casual Bengali with 🌱).
// Use %s placeholder for product name.
var fallbackMessages = []string{
    "%s বেছে নিয়ে তুমি প্রায় ৩০%% বর্জ্য কমাতে সাহায্য করছ 🌱 ছোট্ট পদক্ষেপ, বড় পরিবর্তন!",
    "চমৎকার! %s নেওয়ায় পরিবেশ আরও সবুজ হচ্ছে 🌱",
    "%s কিনে তুমি প্রায় ২৫%% বর্জ্য কমাচ্ছো 🌱 keep going!",
    "প্রকৃতি তোমার পাশে হাসছে 🌱 %s এর মতো sustainable choice নিলে ভবিষ্যৎ উজ্জ্বল হয়।",
    "%s বেছে নিয়ে তুমি পৃথিবীকে একটু হালকা করছ 🌱",
}

// generateMotivationalMessage calls OpenRouter (GPT‑4o) to generate
// a short eco-friendly motivational message in Bengali 🌱.
func (h *ProductHandler) generateMotivationalMessage(product repo.Product, score int) string {
    apiKey := os.Getenv("OPENROUTER_API_KEY")
    if apiKey == "" {
        return randomFallback(product.Name)
    }

    var prompt string
    if score < 50 {
        prompt = fmt.Sprintf(
            "User is considering buying %s by %s. Eco Score: %d (low). "+
                "Write a short, casual and friendly eco‑motivational message in Bengali (max 2 sentences). "+
                "Make sure the message feels natural, not formal — like a friend talking. "+
                "Directly mention the product name in a fun way, so it feels personal. "+
                "Encourage them to try a greener alternative, but keep it supportive and light. "+
                "Also mention a realistic percentage of waste saved or environmental benefit, and vary it each time so it feels fresh. "+
                "Always include an eco emoji 🌱.",
            product.Name, product.BrandName, score,
        )
    } else {
        prompt = fmt.Sprintf(
            "User is buying %s by %s. Eco Score: %d (good). "+
                "Write a short, casual and friendly eco‑motivational message in Bengali (max 2 sentences). "+
                "Make sure the message feels natural, not formal — like a friend talking. "+
                "Directly mention the product name in a fun way, so it feels personal. "+
                "Celebrate their choice and highlight a realistic percentage of waste saved or environmental benefit. "+
                "Vary the style each time — sometimes playful, sometimes poetic, sometimes motivational. "+
                "Always include an eco emoji 🌱.",
            product.Name, product.BrandName, score,
        )
    }

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
    req.Header.Set("HTTP-Referer", "https://yourapp.com")
    req.Header.Set("X-Title", "ecoScanAi")

    client := &http.Client{Timeout: 12 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("OpenRouter API error: %v", err)
        return randomFallback(product.Name)
    }
    defer resp.Body.Close()

    respBody, _ := io.ReadAll(resp.Body)
    log.Printf("OpenRouter status: %d", resp.StatusCode)
    log.Println("OpenRouter raw response:", string(respBody))

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

    return randomFallback(product.Name)
}

// randomFallback returns a random pre-written Bengali message with product name
func randomFallback(productName string) string {
    rand.Seed(time.Now().UnixNano())
    msg := fallbackMessages[rand.Intn(len(fallbackMessages))]
    return fmt.Sprintf(msg, productName)
}
