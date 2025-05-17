package handler

import (
	"github.com/farganamar/evv-service/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) ExternalRouter(r chi.Router) {
	r.Route("/user", func(rc chi.Router) {
		rc.Post("/login", h.UserLogin)
	})

	r.Route("/appointment", func(rc chi.Router) {
		rc.Use(h.AuthMiddleware.Authentication(&middleware.ParamAuth{
			Roles: []string{"CAREGIVER"},
		}))
		rc.Get("/list", h.GetAppointmentList)
	})

}
