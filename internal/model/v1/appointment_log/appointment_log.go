package model

import (
	"encoding/json"

	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null/v5"
)

var (
	ALL_COLUMNS = "id, appointment_id, caregiver_id, log_type, log_data, latitude, longitude, timestamp, notes, created_at"
	TableName   = "appointment_logs"
)

type LogType string

const (
	LogTypeCheckIn  LogType = "CHECK_IN"
	LogTypeCheckOut LogType = "CHECK_OUT"
	LogTypeNote     LogType = "NOTE"
)

type AppointmentLog struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	AppointmentID uuid.UUID       `json:"appointment_id" db:"appointment_id"`
	CaregiverID   uuid.UUID       `json:"caregiver_id" db:"caregiver_id"`
	LogType       null.String     `json:"log_type" db:"log_type"`
	LogData       string          `json:"-" db:"log_data"` // Store as string for DB
	LogDataJSON   json.RawMessage `json:"log_data" db:"-"` // For JSON marshaling
	Latitude      null.Float      `json:"latitude" db:"latitude"`
	Longitude     null.Float      `json:"longitude" db:"longitude"`
	Timestamp     null.Time       `json:"timestamp" db:"timestamp"`
	Notes         null.String     `json:"notes" db:"notes"`
	CreatedAt     null.Time       `json:"created_at" db:"created_at"`
}

func (a *AppointmentLog) PrepareForJSON() error {
	if a.LogData != "" {
		a.LogDataJSON = json.RawMessage(a.LogData)
	}
	return nil
}

func (a *AppointmentLog) PrepareForDB() error {
	if len(a.LogDataJSON) > 0 {
		a.LogData = string(a.LogDataJSON)
	}

	return nil
}

type CheckInOutLogData struct {
	Device string `json:"device"`
	IP     string `json:"ip"`
}

type NoteLogData struct {
	Type string `json:"note"`
}
