package repository

import (
	"context"
	"database/sql"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment_log"
)

func (s *AppointmentLogRepositoryImpl) CreateLog(ctx context.Context, arg model.AppointmentLog, tx *sql.Tx) error {
	query := `INSERT INTO ` + model.TableName + ` (id, appointment_id, caregiver_id, log_type, log_data, latitude, longitude, timestamp, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Execute the query
	var execTx *sql.Tx
	if tx != nil {
		execTx = tx
	}
	_, err := s.Exec(ctx, execTx, query,
		arg.ID,
		arg.AppointmentID,
		arg.CaregiverID,
		arg.LogType,
		arg.LogData,
		arg.Latitude,
		arg.Longitude,
		arg.Timestamp,
		arg.Notes,
	)
	if err != nil {
		return err
	}
	return nil
}
