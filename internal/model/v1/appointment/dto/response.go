package dto

import (
	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null/v5"
)

type GetAppointmentsByUserIdResponse struct {
	AppointmentId string        `json:"appointment_id" swaggertype:"string"`
	StartTime     null.Time     `json:"start_time" swaggertype:"string"`
	EndTime       null.Time     `json:"end_time" swaggertype:"string"`
	Status        null.String   `json:"status" swaggertype:"string"`
	Notes         null.String   `json:"notes" swaggertype:"string"`
	CreatedAt     null.Time     `json:"created_at" swaggertype:"string"`
	UpdatedAt     null.Time     `json:"updated_at" swaggertype:"string"`
	CaregiverId   uuid.NullUUID `json:"caregiver_id" swaggertype:"string"`
	ClientId      uuid.NullUUID `json:"client_id" swaggertype:"string"`
	ClientDetail  ClientDetail  `json:"client_detail"`
}

type ClientDetail struct {
	Name      null.String `json:"name" swaggertype:"string"`
	Phone     null.String `json:"phone" swaggertype:"string"`
	Note      null.String `json:"note" swaggertype:"string"`
	Latitude  null.Float  `json:"latitude" swaggertype:"string"`
	Longitude null.Float  `json:"longitude" swaggertype:"string"`
	Address   null.String `json:"address" swaggertype:"string"`
}
