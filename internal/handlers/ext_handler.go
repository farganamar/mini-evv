package handler

import (
	"net/http"

	"github.com/farganamar/evv-service/transport/http/response"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) ExternalRouter(r chi.Router) {
	r.Route("/external", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			response.WithMessage(w, http.StatusOK, "OK")
		})
	})
}
