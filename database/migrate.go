package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

// RunMigrations executes all pending migrations
func (db *DB) RunMigrations() error {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		return errors.New("MIGRATIONS_DIR is not set")
	}

	// Set the migrations directory
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %v", err)
	}

	sqlDB := stdlib.OpenDB(*db.pool.Config().ConnConfig)
	defer sqlDB.Close()

	// Run migrations
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Info().Msg("Migrations completed successfully")
	return nil
}
