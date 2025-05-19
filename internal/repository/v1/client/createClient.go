package repository

import (
	"context"
	"database/sql"

	model "github.com/farganamar/evv-service/internal/model/v1/client"
)

func (s *ClientRepositoryImpl) CreateClient(ctx context.Context, arg model.Client, tx *sql.Tx) (model.Client, error) {
	query := `INSERT INTO clients (id, name, address, latitude, longitude, phone_number, notes, created_at, updated_at, deleted_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) returning id, name, address, latitude, longitude, phone_number, notes, created_at, updated_at, deleted_at`
	args := []interface{}{
		arg.ID,
		arg.Name,
		arg.Address,
		arg.Latitude,
		arg.Longitude,
		arg.PhoneNumber,
		arg.Notes,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.DeletedAt,
	}

	// Execute the query
	var execTx *sql.Tx
	if tx != nil {
		execTx = tx
	}

	rows, err := s.Query(ctx, execTx, query, args...)
	if err != nil {
		return model.Client{}, err
	}

	defer rows.Close()
	if rows.Next() {
		var client model.Client
		if err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.Address,
			&client.Latitude,
			&client.Longitude,
			&client.PhoneNumber,
			&client.Notes,
			&client.CreatedAt,
			&client.UpdatedAt,
			&client.DeletedAt,
		); err != nil {
			return model.Client{}, err
		}
		return client, nil
	} else {
		return model.Client{}, nil
	}

}
