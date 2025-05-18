package repository

import (
	"context"
	"database/sql"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment"
)

func (q *AppointmentRepositoryImpl) UpdateAppointmentStatus(ctx context.Context, arg model.Appointment, tx *sql.Tx) error {
	query := `UPDATE ` + model.TableName + ` SET status = ? WHERE id = ? AND caregiver_id = ? AND deleted_at IS NULL`

	var execTx *sql.Tx
	if tx != nil {
		execTx = tx
	}

	_, err := q.Exec(ctx, execTx, query, arg.Status, arg.ID, arg.CaregiverID)
	if err != nil {
		return err
	}

	return nil
}
