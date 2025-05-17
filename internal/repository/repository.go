package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/rs/zerolog/log"
	"github.com/zorahealth/user-service/helpers/failure"
	"github.com/zorahealth/user-service/infras"
)

type UserRepoInterface interface {
	UserRepositoryInterface
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type UserRepositoryImpl struct {
	DB *infras.PostgresConn
}

func NewUserRepository(db *infras.PostgresConn) *UserRepositoryImpl {
	s := new(UserRepositoryImpl)
	s.DB = db
	return s
}

func (s *UserRepositoryImpl) exec(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) (pgconn.CommandTag, error) {
	var result pgconn.CommandTag
	var err error

	if tx != nil {
		result, err = tx.Exec(ctx, query, args...)
	} else {
		result, err = s.DB.Write.Exec(ctx, query, args...)
	}

	if err != nil {
		log.Error().Err(err).Msg("[exec] failed exec query")
		return result, failure.InternalError(err)
	}

	return result, err

}

func (s *UserRepositoryImpl) query(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) (pgx.Rows, error) {
	var result pgx.Rows
	var err error

	if tx != nil {
		result, err = tx.Query(ctx, query, args...)
	} else {
		result, err = s.DB.Read.Query(ctx, query, args...)
	}

	if err != nil {
		log.Error().Err(err).Msg("[query] failed exec query")
		return result, failure.InternalError(err)
	}

	return result, nil
}

func (s *UserRepositoryImpl) queryRow(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) pgx.Row {
	if tx != nil {
		return tx.QueryRow(ctx, query, args...)
	}
	return s.DB.Read.QueryRow(ctx, query, args...)
}

func (s *UserRepositoryImpl) copyFrom(ctx context.Context, tx pgx.Tx, tableName string, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	var result int64
	var err error

	result, err = s.DB.Write.CopyFrom(ctx, pgx.Identifier{tableName}, columnNames, rowSrc)

	if err != nil {
		log.Error().Err(err).Msg("[copyFrom] failed to copy data")
		return result, failure.InternalError(err)
	}

	return result, nil
}

func (s *UserRepositoryImpl) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.DB.Write.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		log.Error().Err(err).Msg("[BeginTx] failed to begin transaction")
		return tx, failure.InternalError(err)
	}

	return tx, nil
}
