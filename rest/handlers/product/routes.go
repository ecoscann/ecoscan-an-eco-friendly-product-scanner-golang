package product

import (
	"net/http"

	"ecoscan.com/rest/middlewares"
)

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux, mngr *middlewares.Manager) {
	mux.Handle("GET /api/v1/products/barcode/{barcode}", 
		mngr.Chain(http.HandlerFunc(h.GetProduct)),

	)


	mux.Handle("GET /api/v1/products/search", mngr.Chain(http.HandlerFunc(h.SearchProductsByName)))

	
	

	mux.Handle("POST /api/v1/products/request", 
	mngr.Chain(
		http.HandlerFunc(h.ReqProduct), 
		middlewares.AuthMiddleware,
		),
	)
}
