package main

import (
	"ecoscan/handlers/product"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	product.RegisterRoutes(mux)

	http.ListenAndServe(":2020", mux)
}
