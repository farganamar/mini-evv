package repository

import (
	"context"
	"database/sql"

	"github.com/farganamar/evv-service/helpers/failure"
	"github.com/farganamar/evv-service/infras"
	"github.com/rs/zerolog/log"
)

type RepoInterface interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type RepositoryImpl struct {
	DB *infras.SQLiteConn
}

func NewRepository(db *infras.SQLiteConn) *RepositoryImpl {
	s := new(RepositoryImpl)
	s.DB = db
	return s
}

func (s *RepositoryImpl) Exec(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.ExecContext(ctx, query, args...)
	} else {
		result, err = s.DB.DB.ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.Error().Err(err).Msg("[exec] failed exec query")
		return result, failure.InternalError(err)
	}

	return result, err
}

func (s *RepositoryImpl) Query(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error) {
	var result *sql.Rows
	var err error

	if tx != nil {
		result, err = tx.QueryContext(ctx, query, args...)
	} else {
		result, err = s.DB.DB.QueryContext(ctx, query, args...)
	}

	if err != nil {
		log.Error().Err(err).Msg("[query] failed exec query")
		return result, failure.InternalError(err)
	}

	return result, nil
}

func (s *RepositoryImpl) QueryRow(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) *sql.Row {
	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return s.DB.DB.QueryRowContext(ctx, query, args...)
}

// SQLite doesn't have a built-in COPY command like PostgreSQL
// This is a simplified alternative that executes individual inserts
func (s *RepositoryImpl) BulkInsert(ctx context.Context, tx *sql.Tx, query string, valuesList [][]interface{}) (int64, error) {
	var totalRowsAffected int64

	for _, values := range valuesList {
		var result sql.Result
		var err error

		if tx != nil {
			result, err = tx.ExecContext(ctx, query, values...)
		} else {
			result, err = s.DB.DB.ExecContext(ctx, query, values...)
		}

		if err != nil {
			log.Error().Err(err).Msg("[bulkInsert] failed to insert data")
			return totalRowsAffected, failure.InternalError(err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Error().Err(err).Msg("[bulkInsert] failed to get rows affected")
			return totalRowsAffected, failure.InternalError(err)
		}

		totalRowsAffected += rowsAffected
	}

	return totalRowsAffected, nil
}

func (s *RepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := s.DB.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		log.Error().Err(err).Msg("[BeginTx] failed to begin transaction")
		return tx, failure.InternalError(err)
	}

	return tx, nil
}
