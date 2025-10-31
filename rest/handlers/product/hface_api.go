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

var lowScoreFallbacks = []string{
    "বাহ, %s খেলে সত্যিই রিফ্রেশিং লাগে 🌱\nতবে প্লাস্টিক বোতলটা পরিবেশের জন্য ভালো নয়।\nআপনি যদি ক্যান নিতেন, প্রায় ৩০%% বর্জ্য কমানো যেতো।",
    "%s ব্যবহার করলে মজা আছে 🌱\nকিন্তু এর প্যাকেজিংটা টেকসই নয়।\nআপনি যদি কাচ বা ক্যান বেছে নিতেন, প্রায় ২৫%% সেভ করতে পারতেন।",
    "%s খাওয়া দারুণ লাগে 🌱\nকিন্তু প্লাস্টিক বোতলটা প্রকৃতির ক্ষতি করে।\nআপনি যদি বিকল্প নিতেন, waste reduction দ্বিগুণ হতো।",
}

var goodScoreFallbacks = []string{
    "চমৎকার! %s বেছে নিয়ে আপনি দারুণ কাজ করেছেন 🌱\nএই প্যাকেজিংটা তুলনামূলকভাবে পরিবেশবান্ধব।\nএভাবে প্রায় ৪০%% বর্জ্য কমছে।",
    "%s নেওয়ায় আপনি পরিবেশকে সাহায্য করছেন 🌱\nএটা সত্যিই অনুপ্রেরণাদায়ক একটি সিদ্ধান্ত।\nএভাবে প্রায় ৩৫%% waste সেভ হচ্ছে।",
    "%s কিনে আপনি পৃথিবীকে একটু হালকা করেছেন 🌱\nএটা sustainable choice, ভবিষ্যতের জন্য ভালো।\nএভাবে প্রায় ৪৫%% সেভ হচ্ছে।",
}

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
                "Always end with an eco emoji 🌱." + 
			"demo: Coconut Cookie খেতে অনেক মজা এতে কোকোনাট এর একটা ন্যাচারাল ফ্লেভার আছে তবে Plastic Packaging টা কিন্তু চিন্তা করার বিষয়। এবার কেনাকাটায় একটু greener হোন, Alternatives গুলো চেক করুন better অপশন পেলে প্রায় আপনি 30% plastic waste কমাতে আপনার অবদান রাখতে পারবেন। আসুন সবাই মিলে একটু পরিচ্ছন্ন বাংলাদেশ 🇧🇩 গড়ি। ",
            product.Name, product.BrandName, score, product.Name, product.PackagingMaterial, product.PackagingMaterial,
        )
    } else {
        prompt = fmt.Sprintf(
            "Context: The user is scanning a product (%s by %s) with a good eco-score (%d). "+
                "Task: Write a 3-line, celebratory message in casual Bengali (Banglish style). "+
                "Tone: Respectful 'আপনি', enthusiastic, positive, and reinforcing. "+
                "--- "+
                "Guidelines: "+
                "- **Line 1:** Start with 'চমৎকার!' and the product name (%s). Then, add a *short, positive comment* about its experience, flavor, or usage (e.g., '...এর রিচ ফ্লেভারটা দারুণ!' or '...এটা অনেকেই পছন্দ করে!'). "+
                "- **Line 2:** Specifically praise the *packaging* (e.g., 'আর এর প্যাকেজিংটা বেশ sustainable!'). Then, state the user's direct impact with a **realistic percentage (e.g., 30%%-50%%)** (e.g., 'এটা কিনে আপনি প্রায় XX%% waste কমালেন।'). "+
                "- **Line 3:** Give a positive 'green light' (e.g., 'এটা নিশ্চিন্তে নিতে পারেন।') and a short, patriotic motivational closing (e.g., 'এভাবেই আমাদের পরিবেশ রক্ষায় সাহায্য করুন!'). "+
                "Always end with an eco emoji 🌱." +
			"demo: চমৎকার! Aarong Dairy Chocolate Milk এর রিচ চকলেট এর ফ্লেভার অনেক মজা, অনেকের ই পছন্দ এটা। আর এর প্যাকেজিং অনেক sustainable! এটা কিনলে আপনি প্রায় 40% এর বেশি অপচয় কমালেন। এটা নিশ্চিন্তে কিনতে পারেন। এভাবেই বাংলাদেশ এর পরিবেশ রক্ষায় আপনার অবদান রাখুন।"+
			"write in your own way inspire from the demo dont write the same everytime",
            product.Name, product.BrandName, score, product.Name,
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

func randomScoreAwareFallback(productName string, score int) string {
    rand.Seed(time.Now().UnixNano())
    if score < 50 {
        msg := lowScoreFallbacks[rand.Intn(len(lowScoreFallbacks))]
        return fmt.Sprintf(msg, productName)
    }
    msg := goodScoreFallbacks[rand.Intn(len(goodScoreFallbacks))]
    return fmt.Sprintf(msg, productName)
}
