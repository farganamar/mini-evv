package model

import (
	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null/v5"
)

var (
	ALL_COLUMNS = "id, username, email, phone_number, roles, created_at, updated_at, deleted_at"
	TableName   = "users"
)

type User struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	Username    null.String `json:"username" db:"username"`
	Email       null.String `json:"email" db:"email"`
	PhoneNumber null.String `json:"phone_number" db:"phone_number"`
	Roles       string      `json:"roles" db:"roles"`
	CreatedAt   null.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   null.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt   null.Time   `json:"deleted_at" db:"deleted_at"`
}
