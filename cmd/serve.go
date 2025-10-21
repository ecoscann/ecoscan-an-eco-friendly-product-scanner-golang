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
)

func Serve() {

	cnf := config.GetConfig()
	var sslmode string
	if !cnf.DB.EnableSSLMode {
		sslmode = "disable"
	}

	connectStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s",
		cnf.DB.User, cnf.DB.Password, cnf.DB.Name, sslmode,
	)

	db, err := sqlx.Connect("postgres", connectStr)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
		return
	}
	defer db.Close()

	mngr := middlewares.NewManager()
	mngr.Use(
		middlewares.Logger,
		middlewares.CORS,
	)

	fmt.Println(" Database Connected")

	productHandler := product.NewProductHandler(db)
	UserHandler := user.NewUserHandler(db)

	mux := http.NewServeMux()


	productHandler.RegisterRoutes(mux, mngr)
	UserHandler.RegisterRoutes(mux, mngr)
	
	addr := ":"+ strconv.Itoa(cnf.HttpPort)

		http.ListenAndServe(addr, mux)
}
