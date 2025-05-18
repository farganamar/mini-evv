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
