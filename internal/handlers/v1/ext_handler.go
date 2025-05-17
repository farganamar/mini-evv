package handler

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) ExternalRouter(r chi.Router) {
	r.Route("/user", func(rc chi.Router) {
		rc.Post("/login", h.UserLogin)
	})

}
