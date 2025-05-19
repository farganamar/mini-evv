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

type UpdateAppointmentStatusRequest struct {
	AppointmentId    string         `json:"appointment_id" validate:"required"`
	UserID           string         `json:"user_id" validate:"required" swaggerIgnore:"true"`
	Status           string         `json:"status" validate:"required"`
	Latitude         float64        `json:"latitude" `
	Longitude        float64        `json:"longitude" `
	Note             string         `json:"note"`
	TypeOfNote       string         `json:"type_of_note"`
	VerificationCode string         `json:"verification_code"`
	MetadataDevice   MetadataDevice `json:"metadata_device" `
}

type MetadataDevice struct {
	Device string `json:"device"`
	IP     string `json:"ip"`
}

func (d *UpdateAppointmentStatusRequest) Validate() error {
	validator := helpers.GetValidator()
	if err := validator.Struct(d); err != nil {
		errResponse, err := json.Marshal(helpers.ToErrorResponse(err))
		if err != nil {
			return failure.InternalError(err)
		}

		return errors.New(string(errResponse))
	}

	if d.Status == "IN_PROGRESS" && d.TypeOfNote == "" && d.VerificationCode == "" {
		return errors.New("verification code is required when status is IN_PROGRESS")
	}

	if d.Status == "IN_PROGRESS" && d.Latitude == 0 && d.Longitude == 0 && d.TypeOfNote == "" {
		return errors.New("latitude and longitude are required when status is IN_PROGRESS")
	}

	return nil
}

type SeedAppointmentRequest struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

func (d *SeedAppointmentRequest) Validate() error {
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
