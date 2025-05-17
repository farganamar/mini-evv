package handler

import (
	"net/http"

	"github.com/farganamar/evv-service/helpers/failure"
	"github.com/farganamar/evv-service/internal/model/v1/appointment/dto"
	"github.com/farganamar/evv-service/transport/http/middleware"
	"github.com/farganamar/evv-service/transport/http/response"
)

// GetAppointmentList godoc
// @Summary      Get Appointment List
// @Description  Get Appointment List
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Security  Bearer
// @Success 200 {object} response.Base
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 403 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/evv/appointment/list [get]
func (h *Handler) GetAppointmentList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := middleware.AuthValue{}
	var appointmentListRequest dto.GetAppointmentsByUserIdRequest

	if ctx.Value(middleware.ContextKey) != nil {
		auth = ctx.Value(middleware.ContextKey).(middleware.AuthValue)
		appointmentListRequest.UserId = auth.User.UserID
	}

	if err := appointmentListRequest.Validate(); err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	appointmentListResponse, err := h.AppointmentServiceV1.GetAppointmentsByUserId(ctx, appointmentListRequest)
	if err != nil {
		code := failure.GetCode(err)
		response.WithJSON(w, code, nil, err.Error())
		return
	}

	if len(appointmentListResponse) == 0 {
		response.WithJSON(w, http.StatusOK, nil, "No appointments found")
		return
	}

	response.WithJSON(w, http.StatusOK, appointmentListResponse, "OK")
}
