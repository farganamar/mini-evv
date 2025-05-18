package handler

import (
	"github.com/farganamar/evv-service/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) ExternalRouter(r chi.Router) {
	r.Route("/user", func(rc chi.Router) {
		rc.Post("/login", h.UserLogin)
	})

	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware.Authentication(&middleware.ParamAuth{
			Roles: []string{"CAREGIVER"},
		}))
		r.Route("/appointment", func(r chi.Router) {
			r.Get("/list", h.GetAppointmentList)
			r.Get("/{id}/logs", h.GetAppointmentLogs)
		})

	})

}
