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

// Fallback messages for low-score products (3-line, respectful "আপনি" tone)
var lowScoreFallbacks = []string{
    "বাহ, %s খেলে সত্যিই রিফ্রেশিং লাগে 🌱\nতবে প্লাস্টিক বোতলটা পরিবেশের জন্য ভালো নয়।\nআপনি যদি ক্যান নিতেন, প্রায় ৩০%% বর্জ্য কমানো যেতো।",
    "%s ব্যবহার করলে মজা আছে 🌱\nকিন্তু এর প্যাকেজিংটা টেকসই নয়।\nআপনি যদি কাচ বা ক্যান বেছে নিতেন, প্রায় ২৫%% সেভ করতে পারতেন।",
    "%s খাওয়া দারুণ লাগে 🌱\nকিন্তু প্লাস্টিক বোতলটা প্রকৃতির ক্ষতি করে।\nআপনি যদি বিকল্প নিতেন, waste reduction দ্বিগুণ হতো।",
}

// Fallback messages for good-score products (3-line, respectful "আপনি" tone)
var goodScoreFallbacks = []string{
    "চমৎকার! %s বেছে নিয়ে আপনি দারুণ কাজ করেছেন 🌱\nএই প্যাকেজিংটা তুলনামূলকভাবে পরিবেশবান্ধব।\nএভাবে প্রায় ৪০%% বর্জ্য কমছে।",
    "%s নেওয়ায় আপনি পরিবেশকে সাহায্য করছেন 🌱\nএটা সত্যিই অনুপ্রেরণাদায়ক একটি সিদ্ধান্ত।\nএভাবে প্রায় ৩৫%% waste সেভ হচ্ছে।",
    "%s কিনে আপনি পৃথিবীকে একটু হালকা করেছেন 🌱\nএটা sustainable choice, ভবিষ্যতের জন্য ভালো।\nএভাবে প্রায় ৪৫%% সেভ হচ্ছে।",
}

// generateMotivationalMessage calls OpenRouter (GPT‑4o) to generate
// a 3-line eco-friendly motivational message in Bengali 🌱.
func (h *ProductHandler) generateMotivationalMessage(product repo.Product, score int) string {
    apiKey := os.Getenv("OPENROUTER_API_KEY")
    if apiKey == "" {
        return randomScoreAwareFallback(product.Name, score)
    }

    var prompt string
    if score < 60 {
		 
		prompt = fmt.Sprintf(
			"Context: The user is scanning a product (%s by %s) with a low eco-score (%d). "+
				"Task: Write a 3-line, empathetic, and encouraging message in casual Bengali (Banglish style). "+
				"Tone: Respectful 'আপনি', friendly, light-hearted, and non-judgmental. "+
				"--- "+
				"Guidelines: "+
				"- **Line 1:** Acknowledge their interest and connect with the product's appeal. Start with something like 'বাহ, %s!' or 'Ah, %s!' and mention its nice usage (e.g., '... a quick refresh...'). "+
				"- **Line 2:** Gently give a 'heads-up' about the issue. Use the specific packaging material (%s). Phrase it like, 'Just a heads-up, এর %s packaging-টা environment-এর জন্য একটু heavy.' "+
				"- **Line 3:** Empower them. Suggest a positive alternative (like can/glass) and mention a **realistic impact percentage (e.g., 25%%-40%%)**. Encourage them to 'scroll down' to see the better-scoring options you've found for them. "+
				"Always end with an eco emoji 🌱.",
			product.Name, product.BrandName, score, product.Name, product.PackagingMaterial, product.PackagingMaterial,
		)
	} else {
		
		prompt = fmt.Sprintf(
			"Context: The user is scanning a product (%s by %s) with a good eco-score (%d). "+
				"Task: Write a 3-line, celebratory message in casual Bengali (Banglish style). "+
				"Tone: Respectful 'আপনি', enthusiastic, positive, and reinforcing. "+
				"--- "+
				"Guidelines: "+
				"- **Line 1:** Start with excitement! Confirm their choice with enthusiasm (e.g., 'Yes! %s!' or 'চমৎকার! %s!'). "+
				"- **Line 2:** Directly praise *the user's* choice. Connect their action to the positive outcome (e.g., 'আপনি দারুণ একটা sustainable choice করেছেন!' or 'This is a fantastic eco-friendly pick!'). "+
				"- **Line 3:** Quantify *their* positive impact with a **realistic percentage (e.g., 30%%-50%%)**. Make them feel proud (e.g., 'Your choice just saved over XX%% in waste! Keep up the great work!'). "+
				"Always end with an eco emoji 🌱.",
			product.Name, product.BrandName, score, product.Name, product.Name,
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
        return randomScoreAwareFallback(product.Name, score)
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

    return randomScoreAwareFallback(product.Name, score)
}

// randomScoreAwareFallback picks a 3-line fallback message based on eco score
func randomScoreAwareFallback(productName string, score int) string {
    rand.Seed(time.Now().UnixNano())
    if score < 50 {
        msg := lowScoreFallbacks[rand.Intn(len(lowScoreFallbacks))]
        return fmt.Sprintf(msg, productName)
    }
    msg := goodScoreFallbacks[rand.Intn(len(goodScoreFallbacks))]
    return fmt.Sprintf(msg, productName)
}
