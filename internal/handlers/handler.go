package handler

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) Router(r chi.Router) {
	h.ExternalRouter(r)
}
