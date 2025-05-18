package model

import (
	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null/v5"
)

var (
	ALL_COLUMNS = "id, name, address, latitude, longitude, phone_number, notes, created_at, updated_at, deleted_at"
	TableName   = "clients"
)

type Client struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	Name        null.String `json:"name" db:"name"`
	Address     null.String `json:"address" db:"address"`
	Latitude    null.Float  `json:"latitude" db:"latitude"`
	Longitude   null.Float  `json:"longitude" db:"longitude"`
	PhoneNumber null.String `json:"phone_number" db:"phone_number"`
	Notes       null.String `json:"notes" db:"notes"`
	CreatedAt   null.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   null.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt   null.Time   `json:"deleted_at" db:"deleted_at"`
}
