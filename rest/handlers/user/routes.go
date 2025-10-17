package user

import (
	"net/http"
)

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /api/v1/user", http.HandlerFunc(h.CreateUser))
}
