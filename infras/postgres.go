package infras

import (
	"context"
	"fmt"
	"time"

	"github.com/farganamar/evv-service/configs"
	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxConns          = 10
	minConns          = 2
	maxConnLifetime   = 1 * time.Hour
	maxConnIdleTime   = 30 * time.Minute
	healthCheckPeriod = 1 * time.Minute
)

type PostgresConn struct {
	Read  *pgxpool.Pool
	Write *pgxpool.Pool
}

func ProvidePostgresConn(config *configs.Config) *PostgresConn {
	return &PostgresConn{
		Read:  CreatePostgresReadConn(*config),
		Write: CreatePostgresWriteConn(*config),
	}
}

func CreatePostgresReadConn(config configs.Config) *pgxpool.Pool {
	return CreateDBConnection(
		"read",
		config.DB.Postgres.Read.User,
		config.DB.Postgres.Read.Password,
		config.DB.Postgres.Read.Host,
		config.DB.Postgres.Read.Port,
		config.DB.Postgres.Read.Name,
		config.DB.Postgres.Read.Timezone,
	)
}

func CreatePostgresWriteConn(config configs.Config) *pgxpool.Pool {
	return CreateDBConnection(
		"write",
		config.DB.Postgres.Write.User,
		config.DB.Postgres.Write.Password,
		config.DB.Postgres.Write.Host,
		config.DB.Postgres.Write.Port,
		config.DB.Postgres.Write.Name,
		config.DB.Postgres.Write.Timezone,
	)
}

func CreateDBConnection(name, user, password, host, port, dbName, timezone string) *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse connection string")
	}

	poolConfig.MaxConns = maxConns
	poolConfig.MinConns = minConns
	poolConfig.MaxConnLifetime = maxConnLifetime
	poolConfig.MaxConnIdleTime = maxConnIdleTime
	poolConfig.HealthCheckPeriod = healthCheckPeriod
	poolConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal().Err(err).
			Str("name", name).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Failed connecting to database")
	}

	log.Info().
		Str("name", name).
		Str("host", host).
		Str("port", port).
		Str("dbName", dbName).
		Msg("Connected to database")

	return pool
}

// WithTransaction performs queries with transaction
// func (p *PostgresConn) WithTransaction(ctx context.Context) error {
// 	tx, err := p.Write.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(ctx)

// 	err = fn(tx)
// 	if err != nil {
// 		return err
// 	}

//		return tx.Commit(ctx)
//	}

// func (p *Pool) WithTx(ctx context.Context) (pgx.Tx, error) {
// 	return p.pool.Begin(ctx)
// }
