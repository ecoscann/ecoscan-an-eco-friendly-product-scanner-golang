package user

import (
	"net/http"
)

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /api/v1/auth/register", http.HandlerFunc(h.CreateUser))
}
