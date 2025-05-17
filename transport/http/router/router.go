package router

import (
	handler "github.com/farganamar/evv-service/internal/handlers"
	"github.com/go-chi/chi/v5"
)

// DomainHandlers is a struct that contains all domain-specific handlers.
type DomainHandlers struct {
	Handler handler.Handler
}

// Router is the router struct containing handlers.
type Router struct {
	DomainHandlers DomainHandlers
}

// NewRouter creates a new router.
func NewRouter(domainHandlers DomainHandlers) Router {
	return Router{
		DomainHandlers: domainHandlers,
	}
}

// SetupRoutes sets up all routing for this server.
func (r *Router) SetupRoutes(mux *chi.Mux) {
	mux.Route("/v1", func(rc chi.Router) {
		r.DomainHandlers.Handler.Router(rc)
	})

}
