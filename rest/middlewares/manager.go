package middlewares

import "net/http"

type Middleware func(http.Handler) http.Handler

type Manager struct {
	globalmiddlewares []Middleware
}

func NewManager() *Manager {
	return &Manager{
		globalmiddlewares: []Middleware{},
	}
}

	// to add globalmiddelware to the list
func (m *Manager) Use(middlewares ...Middleware) {
	m.globalmiddlewares = append(m.globalmiddlewares, middlewares...)
}

	// applying all the middlewares
func (m *Manager) Chain(handler http.Handler, routeSpecificMiddlewares ...Middleware) http.Handler {
	
	// specific route middleware
	for i := len(routeSpecificMiddlewares)-1; i >= 0; i--{
		handler = routeSpecificMiddlewares[i](handler)
	}

	// apply global middleware
	for i := len(m.globalmiddlewares) - 1; i>=0 ; i--{
		handler = m.globalmiddlewares[i](handler)
	}

	return handler
}