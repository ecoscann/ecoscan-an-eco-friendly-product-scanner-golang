package main

import (
	"ecoscan.com/rest/handlers/product"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	newMux := product.RegisterRoutes(mux)

	http.ListenAndServe(":2020", newMux)
}
