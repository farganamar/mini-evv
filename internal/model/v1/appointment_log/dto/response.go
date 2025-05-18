package dto

import (
	"encoding/json"

	"github.com/guregu/null/v5"
)

type GetAppointmentLogsResponse struct {
	AppointmentID string          `json:"appointment_id"`
	CaregiverID   string          `json:"caregiver_id"`
	LogType       string          `json:"log_type"`
	LogData       json.RawMessage `json:"log_data"`
	Latitude      null.Float      `json:"latitude"`
	Longitude     null.Float      `json:"longitude"`
	Timestamp     null.Time       `json:"timestamp"`
	Notes         null.String     `json:"notes"`
}
