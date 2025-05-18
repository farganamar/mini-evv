package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	model "github.com/farganamar/evv-service/internal/model/v1/client"
)

func (q *ClientRepositoryImpl) GetClientDetail(ctx context.Context, arg model.Client, tx *sql.Tx) (model.Client, error) {
	res := model.Client{}
	query := `SELECT ` + model.ALL_COLUMNS + ` FROM ` + model.TableName + ` WHERE deleted_at IS NULL`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if !arg.ID.IsNil() {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, arg.ID)
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
	rows, err := q.Query(ctx, execTx, query, args...)
	if err != nil {
		return res, err
	}

	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(
			&res.ID,
			&res.Name,
			&res.Address,
			&res.Latitude,
			&res.Longitude,
			&res.PhoneNumber,
			&res.Notes,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		); err != nil {
			return res, err
		}
	}
	if err := rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}
