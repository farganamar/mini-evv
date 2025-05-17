package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/zorahealth/user-service/internal/service"
)

type UserHandler struct {
	UserService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) UserHandler {
	return UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) Router(r chi.Router) {
	h.ExternalRouter(r)
}
