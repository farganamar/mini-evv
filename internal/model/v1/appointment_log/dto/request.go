package dto

import (
	"encoding/json"
	"errors"

	"github.com/farganamar/evv-service/helpers"
	"github.com/farganamar/evv-service/helpers/failure"
)

type GetAppointmentLogsRequest struct {
	UserId        string `json:"user_id" validate:"required" swaggerIgnore:"true"`
	AppointmentId string `json:"appointment_id" validate:"required" swaggerIgnore:"true"`
}

func (d *GetAppointmentLogsRequest) Validate() error {
	validator := helpers.GetValidator()
	if err := validator.Struct(d); err != nil {
		errResponse, err := json.Marshal(helpers.ToErrorResponse(err))
		if err != nil {
			return failure.InternalError(err)
		}

		return errors.New(string(errResponse))
	}

	return nil
}

type CreateAppointmentLogRequest struct {
	LogType       string `json:"log_type" validate:"required"`
	AppointmentId string `json:"appointment_id" validate:"required"`
	CaregiverId   string `json:"caregiver_id" validate:"required"`
	Latitude      string `json:"latitude" validate:"required"`
	Longitude     string `json:"longitude" validate:"required"`
	Notes         string `json:"notes" validate:"required"`
	LogData       string `json:"log_data" validate:"required"`
}
