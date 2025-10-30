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

// generateMotivationalMessage calls OpenRouter (GPT‑4o) to generate
// a short eco-friendly motivational message in Bengali.
// - If score is low: encourage alternatives, but keep it positive.
// - If score is good: praise the choice and highlight benefits.
// - Always ask the AI to mention a realistic % of waste saved or similar benefit,
//   and vary it each time so it feels fresh.
func (h *ProductHandler) generateMotivationalMessage(product repo.Product, score int) string {
    apiKey := os.Getenv("OPENROUTER_API_KEY")
    if apiKey == "" {
        return "পরিবেশবান্ধব পণ্য বেছে নেওয়া আমাদের পৃথিবীকে রক্ষা করতে সাহায্য করে!"
    }

    var prompt string
    if score < 50 {
        prompt = fmt.Sprintf(
            "একজন ব্যবহারকারী %s (%s) কেনার কথা ভাবছেন। ইকো স্কোর: %d (কম)। "+
                "বাংলায় একটি সংক্ষিপ্ত, সহায়ক পরিবেশবান্ধব বার্তা লিখুন (সর্বোচ্চ ২টি বাক্য)। "+
                "তাদেরকে আরও ভালো প্রভাবের জন্য উচ্চতর ইকো স্কোরের বিকল্প চেষ্টা করতে উৎসাহিত করুন। "+
                "বার্তাটি ইতিবাচক ও অনুপ্রেরণামূলক হোক। "+
                "একটি বাস্তবসম্মত বর্জ্য হ্রাসের শতাংশ বা পরিবেশগত সুবিধা উল্লেখ করুন এবং প্রতিবার ভিন্নভাবে লিখুন যাতে বার্তাটি সতেজ মনে হয়।",
            product.Name, product.BrandName, score,
        )
    } else {
        prompt = fmt.Sprintf(
            "একজন ব্যবহারকারী %s (%s) কিনছেন। ইকো স্কোর: %d (ভালো)। "+
                "বাংলায় একটি সংক্ষিপ্ত, অনুপ্রেরণামূলক পরিবেশবান্ধব বার্তা লিখুন (সর্বোচ্চ ২টি বাক্য)। "+
                "বার্তাটি ব্যক্তিগত ও ইতিবাচক হোক। "+
                "একটি বাস্তবসম্মত বর্জ্য হ্রাসের শতাংশ বা পরিবেশগত সুবিধা উল্লেখ করুন এবং প্রতিবার ভিন্নভাবে লিখুন যাতে বার্তাটি সতেজ মনে হয়।",
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
    req.Header.Set("HTTP-Referer", "https://yourapp.com") // optional
    req.Header.Set("X-Title", "ecoScanAi")                // optional

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("OpenRouter API error: %v", err)
        return "আপনার পছন্দ পরিবেশের জন্য ইতিবাচক প্রভাব ফেলছে!"
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

    return "ধন্যবাদ পরিবেশবান্ধব পণ্য বেছে নেওয়ার জন্য — একসাথে আমরা বর্জ্য কমাতে পারি!"
}
