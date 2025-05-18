package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment"
)

func (q *AppointmentRepositoryImpl) GetAppointmentDetail(ctx context.Context, arg model.Appointment, tx *sql.Tx) (model.Appointment, error) {
	res := model.Appointment{}

	query := `SELECT ` + model.ALL_COLUMNS + ` FROM ` + model.TableName + ` WHERE deleted_at IS NULL`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if !arg.ID.IsNil() {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, arg.ID)
		argIndex++
	}

	if !arg.CaregiverID.IsNil() {
		conditions = append(conditions, fmt.Sprintf("caregiver_id = $%d", argIndex))
		args = append(args, arg.CaregiverID)
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
	if row.Next() {
		if err := row.Scan(
			&res.ID,
			&res.ClientID,
			&res.CaregiverID,
			&res.StartTime,
			&res.EndTime,
			&res.Status,
			&res.VerificationCode,
			&res.Notes,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		); err != nil {
			return res, err
		}
	}

	if err := row.Err(); err != nil {
		return res, err
	}

	return res, nil
}
