package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"ecoscan.com/config"
	"ecoscan.com/rest/handlers/product"
	"ecoscan.com/rest/handlers/user"
	"ecoscan.com/rest/middlewares"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // make sure pq driver is imported
)

func Serve() {
	// 1. Load the config *first*!
	cnf := config.GetConfig()

	// 2. Connect using the DatabaseURL from the config
	db, err := sqlx.Connect("postgres", cnf.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	log.Println("Database Connected")

	// 3. Setup Manager and Handlers (your code was perfect here)
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

	// 4. Listen on the port from the config
	addr := ":" + strconv.Itoa(cnf.HttpPort)

	log.Printf("Server running on %s\n", addr)
	// 5. Use the Manager's chain on the main mux
	http.ListenAndServe(addr, mngr.Chain(mux))
}