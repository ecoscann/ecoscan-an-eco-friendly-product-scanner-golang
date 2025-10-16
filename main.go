package main

import (
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
	"ecoscan.com/rest/handlers/product"
	"github.com/jmoiron/sqlx"
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
	if err != nil{
		fmt.Println("DB Error")
		return
	}

	productHandler := product.NewProductHandler(db)

	mux := http.NewServeMux()

	productHandler.RegisterRoutes(mux)

	http.ListenAndServe(":2020", mux)
}
