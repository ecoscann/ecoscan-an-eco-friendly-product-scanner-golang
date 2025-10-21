package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"ecoscan.com/config"
	"ecoscan.com/rest/handlers/product"
	"ecoscan.com/rest/handlers/user"
	"ecoscan.com/rest/middlewares"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

func Serve() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

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

	addr := ":" + strconv.Itoa(cnf.HttpPort)

	http.ListenAndServe(addr, mux)
}
