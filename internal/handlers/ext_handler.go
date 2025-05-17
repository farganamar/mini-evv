package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zorahealth/user-service/transport/http/response"
)

func (h *UserHandler) ExternalRouter(r chi.Router) {
	r.Route("/external", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			response.WithMessage(w, http.StatusOK, "OK")
		})
	})
}
