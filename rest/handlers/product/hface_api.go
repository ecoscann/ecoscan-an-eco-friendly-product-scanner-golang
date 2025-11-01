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
			"this is an api call"+
		"Context: The user scanned %s by %s. Eco‑score: %d (low).\n"+
		"Task: Write exactly 3 lines in Bengali (Banglish style).\n"+
		"- Line 1: Say something nice about the product.\n"+
		"- Line 2: Point out the environmental issue with its packaging (%s).\n"+
		"- Line 3: Encourage the user to check the alternative products list shown in the app, "+
		"and explain they could reduce waste by choosing one of those greener options.\n"+
		"Tone: Respectful 'আপনি', friendly, motivational, and empowering.\n"+
		"Always end with 🌱.\n\n"+
		"Demo (for inspiration, don’t copy exactly): Coconut Cookie খেতে অনেক মজা... তবে Plastic Packaging টা চিন্তার বিষয়। এবার greener হোন, Alternatives গুলো চেক করুন, better অপশন পেলে প্রায় 30%% plastic waste কমাতে পারবেন। আসুন সবাই মিলে পরিচ্ছন্ন বাংলাদেশ 🇧🇩 গড়ি।",
		product.Name, product.BrandName, score, product.PackagingMaterial,
)

    } else {
        prompt = fmt.Sprintf(
			"this is an api call"+
            "Context: The user scanned %s by %s. Eco‑score: %d (good).\n"+
                "Task: Write exactly 3 lines in Bengali (Banglish style).\n"+
                "- Use respectful 'আপনি' tone.\n"+
                "- Line 1: Mention the product name(in english) and celebrate its taste/usage.\n"+
                "- Line 2: Praise its eco‑friendly packaging or choice.\n"+
                "- Line 3: Highlight a realistic %% waste saved and encourage continuing.\n"+
                "Always end with 🌱.\n\n"+
                "Demo (for inspiration, don’t copy exactly, rewrite in your own way):\n"+
                "চমৎকার! Aarong Dairy Chocolate Milk এর রিচ চকলেট এর ফ্লেভার অনেক মজা, অনেকের ই পছন্দ এটা। আর এর প্যাকেজিং অনেক sustainable! এটা কিনলে আপনি প্রায় 40%% এর বেশি অপচয় কমালেন। এটা নিশ্চিন্তে কিনতে পারেন। এভাবেই বাংলাদেশ এর পরিবেশ রক্ষায় আপনার অবদান রাখুন।",
            product.Name, product.BrandName, score,
        )
    }

    messages := []map[string]string{
        {"role": "user", "content": prompt},
    }

    payload := map[string]interface{}{
        "model": "meta-llama/llama-4-maverick:free", 
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
