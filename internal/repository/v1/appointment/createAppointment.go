package repository

import (
	"context"
	"database/sql"

	model "github.com/farganamar/evv-service/internal/model/v1/appointment"
)

func (s *AppointmentRepositoryImpl) CreateAppointment(ctx context.Context, arg model.Appointment, tx *sql.Tx) (model.Appointment, error) {
	query := `INSERT INTO appointments (id, client_id, caregiver_id, start_time, end_time, status, verification_code, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) returning id, client_id, caregiver_id, start_time, end_time, status, verification_code, created_at, updated_at`

	args := []interface{}{
		arg.ID,
		arg.ClientID,
		arg.CaregiverID,
		arg.StartTime,
		arg.EndTime,
		arg.Status,
		arg.VerificationCode,
		arg.CreatedAt,
		arg.UpdatedAt,
	}

	// Execute the query
	var execTx *sql.Tx
	if tx != nil {
		execTx = tx
	}

	rows, err := s.Query(ctx, execTx, query, args...)
	if err != nil {
		return model.Appointment{}, err
	}
	defer rows.Close()
	if rows.Next() {
		var appointment model.Appointment
		if err := rows.Scan(
			&appointment.ID,
			&appointment.ClientID,
			&appointment.CaregiverID,
			&appointment.StartTime,
			&appointment.EndTime,
			&appointment.Status,
			&appointment.VerificationCode,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		); err != nil {
			return model.Appointment{}, err
		}
		return appointment, nil
	} else {
		return model.Appointment{}, nil
	}

}
