package handler

import (
	UserServiceV1 "github.com/farganamar/evv-service/internal/service/v1/user"
	"github.com/farganamar/evv-service/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	AuthMiddleware middleware.AuthMiddlewareInterface
	UserServiceV1  UserServiceV1.UserService
}

func NewHandler(
	authMiddleware middleware.AuthMiddlewareInterface,
	userServiceV1 UserServiceV1.UserService,
) Handler {
	return Handler{
		AuthMiddleware: authMiddleware,
		UserServiceV1:  userServiceV1,
	}
}

func (h *Handler) Router(r chi.Router) {
	r.Route("/evv", func(r chi.Router) {
		h.ExternalRouter(r)
	})
}
