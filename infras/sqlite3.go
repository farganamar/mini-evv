package infras

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/farganamar/evv-service/configs"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/rs/zerolog/log"
)

// SQLiteConn represents a SQLite database connection
type SQLiteConn struct {
	DB *sql.DB
}

// ProvideSQLiteConn creates and returns a new SQLite connection
func ProvideSQLiteConn(config *configs.Config) *SQLiteConn {
	return &SQLiteConn{
		DB: CreateSQLiteConnection(config.DB.SQLite.Path),
	}
}

// CreateSQLiteConnection establishes a connection to the SQLite database
func CreateSQLiteConnection(dbPath string) *sql.DB {
	// Ensure the directory exists
	dir := getDirectoryFromPath(dbPath)
	if dir != "" && dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal().Err(err).Str("path", dir).Msg("Failed to create directory for SQLite database")
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", dbPath).Msg("Failed to open SQLite database")
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Str("path", dbPath).Msg("Failed to connect to SQLite database")
	}

	log.Info().Str("path", dbPath).Msg("Connected to SQLite database")
	return db
}

// Close closes the SQLite database connection
func (s *SQLiteConn) Close() error {
	if s.DB != nil {
		err := s.DB.Close()
		if err != nil {
			return fmt.Errorf("error closing SQLite connection: %w", err)
		}
		log.Info().Msg("SQLite connection closed")
	}
	return nil
}

// WithTransaction executes a function within a transaction
func (s *SQLiteConn) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Helper function to extract directory from path
func getDirectoryFromPath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return ""
}
