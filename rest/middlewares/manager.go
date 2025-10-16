package middlewares

import "net/http"

type Middleware func(http.Handler) http.Handler

type Manager struct {
	globalmiddlewares []Middleware
}

func NewManager() *Manager {
	return &Manager{
		globalmiddlewares: make([]Middleware, 0),
	}
}

func (mnger *Manager) Use(middlewares ...Middleware) {
	mnger.globalmiddlewares = append(mnger.globalmiddlewares, middlewares...)
}

func (mnger *Manager) With(handler http.Handler, middlewares ...Middleware) http.Handler {
	hd := handler
	for _, middleware := range middlewares {
		hd = middleware(hd)
	}
	return hd
}

func (mnger *Manager) WrapMux(handler http.Handler) http.Handler {
	hd := handler
	for _, middleware := range mnger.globalmiddlewares {
		hd = middleware(hd)
	}
	return hd
}
