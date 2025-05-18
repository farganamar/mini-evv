package handler

import (
	"net/http"

	"github.com/farganamar/evv-service/helpers"
	"github.com/farganamar/evv-service/helpers/failure"
	"github.com/farganamar/evv-service/internal/model/v1/appointment/dto"
	AppointmentLogDTO "github.com/farganamar/evv-service/internal/model/v1/appointment_log/dto"
	"github.com/farganamar/evv-service/transport/http/middleware"
	"github.com/farganamar/evv-service/transport/http/response"
	"github.com/go-chi/chi/v5"
)

// GetAppointmentList godoc
// @Summary      Get Appointment List
// @Description  Get Appointment List
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param status query string false "status" Enums(SCHEDULED, CANCELED, COMPLETED, IN_PROGRESS)
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

	if err := helpers.ParseQueryParams(r, &appointmentListRequest); err != nil {
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

// GetAppointmentLogs godoc
// @Summary      Get Appointment Logs
// @Description  Get Appointment Logs
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param id path string true "Appointment ID"
// @Security  Bearer
// @Success 200 {object} response.Base
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 403 {object} response.Base
// @Failure 404 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/evv/appointment/{id}/logs [get]
func (h *Handler) GetAppointmentLogs(w http.ResponseWriter, r *http.Request) {
	appointmentId := chi.URLParam(r, "id")
	ctx := r.Context()
	auth := middleware.AuthValue{}
	var appointmentLogsRequest AppointmentLogDTO.GetAppointmentLogsRequest
	appointmentLogsRequest.AppointmentId = appointmentId

	if ctx.Value(middleware.ContextKey) != nil {
		auth = ctx.Value(middleware.ContextKey).(middleware.AuthValue)
		appointmentLogsRequest.UserId = auth.User.UserID
	}

	if err := appointmentLogsRequest.Validate(); err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	appointmentLogsResponse, err := h.AppointmentLogServiceV1.GetAppointmentLogs(ctx, appointmentLogsRequest)
	if err != nil {
		code := failure.GetCode(err)
		response.WithJSON(w, code, nil, err.Error())
		return
	}

	if len(appointmentLogsResponse) == 0 {
		response.WithJSON(w, http.StatusOK, nil, "No appointment logs found")
		return
	}

	response.WithJSON(w, http.StatusOK, appointmentLogsResponse, "OK")
}
