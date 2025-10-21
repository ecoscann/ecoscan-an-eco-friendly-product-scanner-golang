package cmd

import (
	"log"
	"net/http"
	"strconv"

	"ecoscan.com/config"
	"ecoscan.com/rest/handlers/product"
	"ecoscan.com/rest/handlers/user"
	"ecoscan.com/rest/middlewares"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Serve() {
	// ১. কনফিগ লোড করুন
	cnf := config.GetConfig()

	// ২. DATABASE_URL দিয়ে সরাসরি কানেক্ট করুন
	db, err := sqlx.Connect("postgres", cnf.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	log.Println("Database Connected")

	// ৩. আপনার বাকি কোড (ম্যানেজার, হ্যান্ডলার) ঠিকই আছে
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