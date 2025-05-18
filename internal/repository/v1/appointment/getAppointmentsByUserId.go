package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment"
)

func (q *AppointmentRepositoryImpl) GetAppointmentsByUserId(ctx context.Context, arg model.Appointment, tx *sql.Tx) ([]model.Appointment, error) {
	res := []model.Appointment{}

	query := `SELECT ` + model.ALL_COLUMNS + ` FROM ` + model.TableName + ` WHERE deleted_at IS NULL`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if !arg.CaregiverID.IsNil() {
		conditions = append(conditions, fmt.Sprintf("caregiver_id = $%d", argIndex))
		args = append(args, arg.CaregiverID)
		argIndex++
	}

	if arg.Status.Valid && arg.Status.String != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, arg.Status.String)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	} else {

		return res, nil
	}

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
		var appointment model.Appointment
		if err := row.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.CaregiverID,
			&appointment.StartTime,
			&appointment.EndTime,
			&appointment.Status,
			&appointment.VerificationCode,
			&appointment.Notes,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
			&appointment.DeletedAt,
		); err != nil {
			return res, err
		}

		res = append(res, appointment)
	}

	return res, nil
}
