package cmd

import (
	"fmt"
	"net/http"

	"ecoscan.com/config"
	"ecoscan.com/rest/handlers/product"
	"ecoscan.com/rest/handlers/user"
	"github.com/jmoiron/sqlx"
)

func Serve() {

	cnf := config.GetConfig()

	sslmode := "disable"
	if cnf.DB.EnableSSLMode {
		sslmode = "Required"
	}

	connectStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s",
		cnf.DB.User, cnf.DB.Password,cnf.DB.Name, sslmode
	)

	db, err := sqlx.Connect("postgres", connectStr)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
		return
	}
	defer db.Close()

	fmt.Println(" Database Connected")

	productHandler := product.NewProductHandler(db)
	UserHandler := user.NewUserHandler(db)

	mux := http.NewServeMux()

	productHandler.RegisterRoutes(mux)
	UserHandler.RegisterRoutes(mux)
	fmt.Sprintf(": %d")
	addr:=

	http.ListenAndServe(":2020", mux)
}
