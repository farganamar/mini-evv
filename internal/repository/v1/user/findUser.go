package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	model "github.com/farganamar/evv-service/internal/model/v1/user"
)

func (q *UserRepositoryImpl) FindUser(ctx context.Context, arg model.User, tx *sql.Tx) (model.User, error) {
	res := model.User{}

	query := `SELECT ` + model.ALL_COLUMNS + ` FROM ` + model.TableName + ` WHERE deleted_at IS NULL`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if !arg.ID.IsNil() {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, arg.ID)
		argIndex++
	}

	if arg.Username.Valid && arg.Username.String != "" {
		conditions = append(conditions, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, arg.Username)
		argIndex++
	}

	if arg.Email.Valid && arg.Email.String != "" {
		conditions = append(conditions, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, arg.Email)
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

	row := q.QueryRow(ctx, execTx, query, args...)

	if err := row.Scan(
		&res.ID,
		&res.Username,
		&res.Email,
		&res.PhoneNumber,
		&res.Roles,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("failed to scan user: %w", err)
	}

	return res, nil

}
