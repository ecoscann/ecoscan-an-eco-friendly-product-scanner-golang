package user

import (
	"net/http"

	"ecoscan.com/rest/middlewares"
)

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux, mngr *middlewares.Manager) {

	mux.Handle("POST /api/v1/auth/register",
	mngr.Chain(http.HandlerFunc(h.CreateUser), 
		),	
	)

	mux.Handle("POST /api/v1/auth/login",
	mngr.Chain(http.HandlerFunc(h.LoginUser),
		),
	)
}
