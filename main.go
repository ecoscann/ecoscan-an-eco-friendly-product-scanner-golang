package main

import (
	"fmt"
	"net/http"

	"ecoscan.com/rest/handlers/product"
	"ecoscan.com/rest/handlers/user"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

/* DB_HOST=localhost
DB_PORT=5432
DB_NAME=ecommerce
DB_USER=postgres
DB_PASSWORD=1212
DB_ENABLE_SSL_MODE=false */

func main() {

	connectStr := "user=postgres password=1212 dbname=ecoscan sslmode=disable"

	db, err := sqlx.Connect("postgres", connectStr)
	if err != nil {
		fmt.Println("DB Error")
		return
	}

	productHandler := product.NewProductHandler(db)
	UserHandler := user.NewUserHandler(db)

	mux := http.NewServeMux()

	productHandler.RegisterRoutes(mux)
	UserHandler.RegisterRoutes(mux)

	http.ListenAndServe(":2020", mux)
}
