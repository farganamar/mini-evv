package router

import (
	"net/http"

	handlerV1 "github.com/farganamar/evv-service/internal/handlers/v1"
	"github.com/farganamar/evv-service/transport/http/response"
	"github.com/go-chi/chi/v5"
)

// DomainHandlers is a struct that contains all domain-specific handlers.
type DomainHandlers struct {
	HandlerV1 handlerV1.Handler
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
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.WithMessage(w, http.StatusOK, "OK")
	})

	mux.Route("/v1", func(rc chi.Router) {
		r.DomainHandlers.HandlerV1.Router(rc)
	})

}
