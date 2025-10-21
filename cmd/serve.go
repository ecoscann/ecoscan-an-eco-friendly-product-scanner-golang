package cmd

import (
    "log"
    "net/http"
    "os"
    "strconv"

    "ecoscan.com/config"
    "ecoscan.com/rest/handlers/product"
    "ecoscan.com/rest/handlers/user"
    "ecoscan.com/rest/middlewares"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq" // make sure pq driver is imported
)

func Serve() {
    // connect to Supabase using DATABASE_URL
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL is not set")
    }

    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("Database connection error: %v", err)
    }
    defer db.Close()

    mngr := middlewares.NewManager()
    mngr.Use(
        middlewares.Logger,
        middlewares.CORS,
    )

    log.Println("Database Connected")

    productHandler := product.NewProductHandler(db)
    userHandler := user.NewUserHandler(db)

    mux := http.NewServeMux()
    productHandler.RegisterRoutes(mux, mngr)
    userHandler.RegisterRoutes(mux, mngr)

    cnf := config.GetConfig()
    addr := ":" + strconv.Itoa(cnf.HttpPort)

    log.Printf("Server running on %s\n", addr)
    http.ListenAndServe(addr, mux)
}
