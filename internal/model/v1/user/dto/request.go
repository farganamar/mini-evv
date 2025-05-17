package dto

import (
	"encoding/json"
	"errors"

	"github.com/farganamar/evv-service/helpers"
	"github.com/farganamar/evv-service/helpers/failure"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required" swaggertype:"string"`
}

func (d *LoginRequest) Validate() error {
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
