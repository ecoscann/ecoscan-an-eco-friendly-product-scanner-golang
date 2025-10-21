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
	
	cnf := config.GetConfig()

	
	db, err := sqlx.Connect("postgres", cnf.DatabaseURL)
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