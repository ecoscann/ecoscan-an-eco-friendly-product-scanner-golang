package cmd

import (
	"fmt" // <-- 'fmt' ইম্পোর্ট করুন
	"log"
	"net" // <-- 'net' প্যাকেজ ইম্পোর্ট করুন
	"net/http"
	"strconv"

	"ecoscan.com/config"
	"ecoscan.com/rest/handlers/product"
	"ecoscan.com/rest/handlers/user"
	"ecoscan.com/rest/middlewares"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// resolveHostToIPv4 হোস্টনেমকে তার প্রথম IPv4 অ্যাড্রেসে পরিণত করে
func resolveHostToIPv4(host string) (string, error) {
	// "localhost" পরিবর্তন করবেন না
	if host == "localhost" {
		return "127.0.0.1", nil
	}
	
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil // প্রথম যে IPv4 অ্যাড্রেসটি পাবে, সেটি রিটার্ন করবে
		}
	}
	return "", fmt.Errorf("no IPv4 address found for %s", host)
}

func Serve() {
	cnf := config.GetConfig()

	// --- IPv4 ফিক্স শুরু ---
	// Render-এ ডেপ্লয় করার জন্য হোস্টনেমকে IPv4-তে পরিণত করুন
	dbHost, err := resolveHostToIPv4(cnf.DB.Host)
	if err != nil {
		log.Fatalf("Could not resolve DB host to IPv4: %v", err)
	}
	// --- IPv4 ফিক্স শেষ ---

	var sslmode string
	if cnf.DB.EnableSSLMode {
		sslmode = "require" // Supabase-এর জন্য "require"
	} else {
		sslmode = "disable" // লোকালহোস্টের জন্য "disable"
	}

	// IPv4 অ্যাড্রেস দিয়ে কানেকশন স্ট্রিং তৈরি করুন
	connectStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, // এখানে cnf.DB.Host-এর বদলে dbHost ব্যবহার করুন
		cnf.DB.Port,
		cnf.DB.User,
		cnf.DB.Password,
		cnf.DB.Name,
		sslmode,
	)

	db, err := sqlx.Connect("postgres", connectStr)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	log.Println("Database Connected")

	mngr := middlewares.NewManager()
	mngr.Use(
		middlewares.Logger,
		middlewares.CORS,
	)

	productHandler := product.NewProductHandler(db)
	userHandler := user.NewUserHandler(db)

	mux := http.NewServeMux()
	productHandler.RegisterRoutes(mux, mngr)
	userHandler.RegisterRoutes(mux, mngr)

	addr := ":" + strconv.Itoa(cnf.HttpPort)

	log.Printf("Server running on %s\n", addr)
	http.ListenAndServe(addr, mngr.Chain(mux))
}