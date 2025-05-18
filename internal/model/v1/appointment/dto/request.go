package dto

import (
	"encoding/json"
	"errors"

	"github.com/farganamar/evv-service/helpers"
	"github.com/farganamar/evv-service/helpers/failure"
)

type GetAppointmentsByUserIdRequest struct {
	UserId string `json:"user_id" validate:"required" swaggerIgnore:"true"`
	Status string `query:"status" `
}

func (d *GetAppointmentsByUserIdRequest) Validate() error {
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
