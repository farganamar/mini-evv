package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
	model "github.com/zorahealth/user-service/internal/model/user"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    first_name,
    last_name,
    phone_number,
    phone_country_code,
    password,
    profile_picture,
    is_verified,
    created_at,
    updated_at,
    deleted_at,
    dob,
    gender
) VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5, 
    $6, 
    $7, 
    $8, 
    $9, 
    $10,
    $11,
    $12,
    $13,
    $14 
)
RETURNING id, email, first_name, last_name, phone_number, phone_country_code, password, profile_picture, is_verified, created_at, updated_at, deleted_at, dob, gender
`

func (q *UserRepositoryImpl) CreateUser(ctx context.Context, arg *model.User, tx pgx.Tx) (model.User, error) {
	row := q.queryRow(ctx, tx,
		createUser,
		arg.ID,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.PhoneNumber,
		arg.PhoneCountryCode,
		arg.Password,
		arg.ProfilePicture,
		arg.IsVerified,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.DeletedAt,
		arg.Dob,
		arg.Gender,
	)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.PhoneNumber,
		&i.PhoneCountryCode,
		&i.Password,
		&i.ProfilePicture,
		&i.IsVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Dob,
		&i.Gender,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *UserRepositoryImpl) DeleteUser(ctx context.Context, id pgtype.UUID, tx pgx.Tx) error {
	_, err := q.exec(ctx, tx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, email, first_name, last_name, phone_number, phone_country_code, password, profile_picture, is_verified, created_at, updated_at, deleted_at, dob, gender FROM users
WHERE id = $1 LIMIT 1
`

func (q *UserRepositoryImpl) GetUser(ctx context.Context, id uuid.UUID, tx pgx.Tx) (model.User, error) {
	row := q.queryRow(ctx, tx, getUser, id)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.PhoneNumber,
		&i.PhoneCountryCode,
		&i.Password,
		&i.ProfilePicture,
		&i.IsVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Dob,
		&i.Gender,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, first_name, last_name, phone_number, phone_country_code, password, profile_picture, is_verified, created_at, updated_at, deleted_at, dob, gender FROM users
ORDER BY id
`

func (q *UserRepositoryImpl) ListUsers(ctx context.Context, tx pgx.Tx) ([]model.User, error) {
	rows, err := q.query(ctx, tx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.User
	for rows.Next() {
		var i model.User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.FirstName,
			&i.LastName,
			&i.PhoneNumber,
			&i.PhoneCountryCode,
			&i.Password,
			&i.ProfilePicture,
			&i.IsVerified,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Dob,
			&i.Gender,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *UserRepositoryImpl) UpdateUser(ctx context.Context, arg model.User, tx pgx.Tx) error {
	var setClauses []string
	var args []interface{}

	args = append(args, arg.ID)
	argIndex := 2

	if arg.FirstName != nil {
		setClauses = append(setClauses, fmt.Sprintf("first_name = $%d", argIndex))
		args = append(args, *arg.FirstName)
		argIndex++
	}

	if arg.LastName != nil {
		setClauses = append(setClauses, fmt.Sprintf("last_name = $%d", argIndex))
		args = append(args, *arg.LastName)
		argIndex++
	}

	if arg.PhoneNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone_number = $%d", argIndex))
		args = append(args, *arg.PhoneNumber)
		argIndex++
	}

	if arg.PhoneCountryCode != "" {
		setClauses = append(setClauses, fmt.Sprintf("phone_country_code = $%d", argIndex))
		args = append(args, arg.PhoneCountryCode)
		argIndex++
	}

	if arg.Password != nil {
		setClauses = append(setClauses, fmt.Sprintf("password = $%d", argIndex))
		args = append(args, *arg.Password)
		argIndex++
	}

	if arg.ProfilePicture != nil {
		setClauses = append(setClauses, fmt.Sprintf("profile_picture = $%d", argIndex))
		args = append(args, *arg.ProfilePicture)
		argIndex++
	}

	if arg.IsVerified != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_verified = $%d", argIndex))
		args = append(args, *arg.IsVerified)
		argIndex++
	}

	if arg.DeletedAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("deleted_at = $%d", argIndex))
		args = append(args, *arg.DeletedAt)
		argIndex++
	}

	if !arg.Dob.IsZero() {
		setClauses = append(setClauses, fmt.Sprintf("dob = $%d", argIndex))
		args = append(args, &arg.Dob)
		argIndex++
	}

	if arg.VerifyAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("verify_at = $%d", argIndex))
		args = append(args, *arg.VerifyAt)
		argIndex++
	}

	if len(setClauses) == 0 {
		return nil
	}

	updatedAt := time.Now()
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, updatedAt)
	argIndex++

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $1", strings.Join(setClauses, ", "))

	_, err := q.exec(ctx, tx, query, args...)
	return err
}

func (q *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string, tx pgx.Tx) (model.User, error) {
	query := `SELECT id, email, first_name, last_name, phone_number, phone_country_code, password, profile_picture, is_verified, created_at, updated_at, deleted_at, dob, gender, verify_at FROM users
WHERE email = $1 AND deleted_at ISNULL LIMIT 1`
	row := q.queryRow(ctx, tx, query, email)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.PhoneNumber,
		&i.PhoneCountryCode,
		&i.Password,
		&i.ProfilePicture,
		&i.IsVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Dob,
		&i.Gender,
		&i.VerifyAt,
	)
	return i, err
}

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, arg *model.User, tx pgx.Tx) (model.User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID, tx pgx.Tx) error
	GetUser(ctx context.Context, id uuid.UUID, tx pgx.Tx) (model.User, error)
	ListUsers(ctx context.Context, tx pgx.Tx) ([]model.User, error)
	UpdateUser(ctx context.Context, arg model.User, tx pgx.Tx) error
	GetUserByEmail(ctx context.Context, email string, tx pgx.Tx) (model.User, error)
}
