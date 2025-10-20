package product

import (
	"net/http"
)

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/products/barcode/{barcode}", http.HandlerFunc(h.GetProduct))
	mux.Handle("POST /api/v1/request", http.HandlerFunc(h.ReqProduct))
}
