package handler

import (
	AppointmentServiceV1 "github.com/farganamar/evv-service/internal/service/v1/appointment"
	UserServiceV1 "github.com/farganamar/evv-service/internal/service/v1/user"
	"github.com/farganamar/evv-service/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	AuthMiddleware       middleware.AuthMiddlewareInterface
	UserServiceV1        UserServiceV1.UserService
	AppointmentServiceV1 AppointmentServiceV1.AppointmentService
}

func NewHandler(
	authMiddleware middleware.AuthMiddlewareInterface,
	userServiceV1 UserServiceV1.UserService,
	appointmentServiceV1 AppointmentServiceV1.AppointmentService,
) Handler {
	return Handler{
		AuthMiddleware:       authMiddleware,
		UserServiceV1:        userServiceV1,
		AppointmentServiceV1: appointmentServiceV1,
	}
}

func (h *Handler) Router(r chi.Router) {
	r.Route("/evv", func(r chi.Router) {
		h.ExternalRouter(r)
	})
}
