package model

import (
	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null/v5"
)

var (
	ALL_COLUMNS = "id, client_id, caregiver_id, start_time, end_time, status, verification_code, notes, created_at, updated_at, deleted_at"
	TableName   = "appointments"
)

const (
	StatusShecheduled = "SCHEDULED"
	StatusCompleted   = "COMPLETED"
	StatusCancelled   = "CANCELLED"
	StatusInProgress  = "IN_PROGRESS"
)

type Appointment struct {
	ID               uuid.UUID   `json:"id" db:"id"`
	ClientID         uuid.UUID   `json:"client_id" db:"client_id"`
	CaregiverID      uuid.UUID   `json:"caregiver_id" db:"caregiver_id"`
	StartTime        null.Time   `json:"start_time" db:"start_time"`
	EndTime          null.Time   `json:"end_time" db:"end_time"`
	Status           null.String `json:"status" db:"status"`
	VerificationCode null.String `json:"verification_code" db:"verification_code"`
	Notes            null.String `json:"notes" db:"notes"`
	CreatedAt        null.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        null.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt        null.Time   `json:"deleted_at" db:"deleted_at"`
}
