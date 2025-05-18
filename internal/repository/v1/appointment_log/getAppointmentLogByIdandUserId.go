package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment_log"
)

func (q *AppointmentLogRepositoryImpl) GetAppointmentLogByIdandUserId(ctx context.Context, arg model.AppointmentLog, tx *sql.Tx) ([]model.AppointmentLog, error) {
	res := []model.AppointmentLog{}

	query := `SELECT ` + model.ALL_COLUMNS + ` FROM ` + model.TableName

	var conditions []string
	var args []interface{}
	argIndex := 1

	if !arg.CaregiverID.IsNil() {
		conditions = append(conditions, fmt.Sprintf("caregiver_id = $%d", argIndex))
		args = append(args, arg.CaregiverID)
		argIndex++
	}

	if !arg.AppointmentID.IsNil() {
		conditions = append(conditions, fmt.Sprintf("appointment_id = $%d", argIndex))
		args = append(args, arg.AppointmentID)
		argIndex++
	}

	if arg.LogType.Valid && arg.LogType.String != "" {
		conditions = append(conditions, fmt.Sprintf("log_type = $%d", argIndex))
		args = append(args, arg.LogType)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	} else {
		return res, nil
	}

	query += " ORDER BY timestamp ASC"

	// Execute the query
	var execTx *sql.Tx
	if tx != nil {
		execTx = tx
	}

	row, err := q.Query(ctx, execTx, query, args...)
	if err != nil {
		return res, err
	}

	defer row.Close()
	for row.Next() {
		var appointmentLog model.AppointmentLog
		if err := row.Scan(
			&appointmentLog.ID,
			&appointmentLog.AppointmentID,
			&appointmentLog.CaregiverID,
			&appointmentLog.LogType,
			&appointmentLog.LogData,
			&appointmentLog.Latitude,
			&appointmentLog.Longitude,
			&appointmentLog.Timestamp,
			&appointmentLog.Notes,
			&appointmentLog.CreatedAt,
		); err != nil {
			return res, err
		}

		if err := appointmentLog.PrepareForJSON(); err != nil {
			return res, err
		}

		res = append(res, appointmentLog)
	}

	return res, nil
}
