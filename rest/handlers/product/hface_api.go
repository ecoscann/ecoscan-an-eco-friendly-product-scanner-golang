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

// Pre-written fallback messages (natural Bengali, casual, with 🌱)
var fallbackMessages = []string{
    "তুমি এই পণ্যটি বেছে নিয়ে প্রায় ৩০% বর্জ্য কমাতে সাহায্য করছ 🌱 ছোট্ট পদক্ষেপ, বড় পরিবর্তন!",
    "চমৎকার! এই সিদ্ধান্তে পরিবেশ আরও সবুজ হচ্ছে 🌱",
    "তোমার এই চয়েসে প্রায় ২৫% বর্জ্য কমছে 🌱 keep going!",
    "প্রকৃতি তোমার পাশে হাসছে 🌱 sustainable choice নিলে ভবিষ্যৎ উজ্জ্বল হয়।",
    "এই পণ্যটি বেছে নিয়ে তুমি পৃথিবীকে একটু হালকা করছ 🌱",
}

// generateMotivationalMessage calls OpenRouter (GPT‑4o) to generate
// a short eco-friendly motivational message in Bengali 🌱.
// - If score is low: encourage alternatives, but keep it supportive.
// - If score is good: praise the choice and highlight benefits.
// - Always in natural, inspiring Bengali (not overly formal).
// - Randomize style: sometimes poetic, sometimes playful, sometimes motivational.
// - Always include an eco emoji 🌱.
// - If API fails, return a random fallback message.
func (h *ProductHandler) generateMotivationalMessage(product repo.Product, score int) string {
    apiKey := os.Getenv("OPENROUTER_API_KEY")
    if apiKey == "" {
        return randomFallback()
    }

    var prompt string
    if score < 50 {
        prompt = fmt.Sprintf(
            "একজন ব্যবহারকারী %s (%s) কেনার কথা ভাবছেন। ইকো স্কোর: %d (কম)। "+
                "বাংলায় একটি সংক্ষিপ্ত, স্বাভাবিক ও অনুপ্রেরণামূলক পরিবেশবান্ধব বার্তা লিখুন (সর্বোচ্চ ২টি বাক্য)। "+
                "তাদেরকে আরও ভালো প্রভাবের জন্য উচ্চতর ইকো স্কোরের বিকল্প চেষ্টা করতে উৎসাহিত করুন। "+
                "বার্তাটি যেন বন্ধুসুলভ ও ইতিবাচক হয়। "+
                "একটি বাস্তবসম্মত বর্জ্য হ্রাসের শতাংশ বা পরিবেশগত সুবিধা উল্লেখ করুন এবং প্রতিবার ভিন্নভাবে লিখুন যাতে বার্তাটি সতেজ মনে হয়। "+
                "স্টাইলটি প্রতিবার ভিন্ন হোক — কখনও কাব্যিক, কখনও খেলাচ্ছলে, কখনও সরাসরি অনুপ্রেরণামূলক। "+
                "বার্তায় একটি পরিবেশ ইমোজি 🌱 ব্যবহার করুন।",
            product.Name, product.BrandName, score,
        )
    } else {
        prompt = fmt.Sprintf(
            "একজন ব্যবহারকারী %s (%s) কিনছেন। ইকো স্কোর: %d (ভালো)। "+
                "বাংলায় একটি সংক্ষিপ্ত, স্বাভাবিক ও অনুপ্রেরণামূলক পরিবেশবান্ধব বার্তা লিখুন (সর্বোচ্চ ২টি বাক্য)। "+
                "বার্তাটি যেন ব্যক্তিগত, উষ্ণ ও ইতিবাচক হয়। "+
                "একটি বাস্তবসম্মত বর্জ্য হ্রাসের শতাংশ বা পরিবেশগত সুবিধা উল্লেখ করুন এবং প্রতিবার ভিন্নভাবে লিখুন যাতে বার্তাটি সতেজ মনে হয়। "+
                "স্টাইলটি প্রতিবার ভিন্ন হোক — কখনও কাব্যিক, কখনও খেলাচ্ছলে, কখনও সরাসরি অনুপ্রেরণামূলক। "+
                "বার্তায় একটি পরিবেশ ইমোজি 🌱 ব্যবহার করুন।",
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

    client := &http.Client{Timeout: 15 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("OpenRouter API error: %v", err)
        return randomFallback()
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

    return randomFallback()
}

// randomFallback returns a random pre-written Bengali message
func randomFallback() string {
    rand.Seed(time.Now().UnixNano())
    return fallbackMessages[rand.Intn(len(fallbackMessages))]
}
