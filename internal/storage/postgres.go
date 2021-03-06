package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDataSource struct {
	DB *sqlx.DB
}

// InitPostgres establishes connections to fields in PostgresDataSource
func InitPostgres() (*PostgresDataSource, error) {
	log.Printf("Initializing data sources\n")

	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", pgHost, pgPort, pgUser, pgPassword, pgDB)

	log.Printf("Connecting to Postgresql\n")
	db, err := sqlx.Open("postgres", pgConnString)

	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Verify database connection is working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("error create driver for migrate: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/storage/migrations/pg",
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("error create migrating: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	log.Printf("Postgres Migration done\n")

	return &PostgresDataSource{
		DB: db,
	}, nil
}

// close to be used in graceful server shutdown
func (d *PostgresDataSource) close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error closing Postgresql: %w", err)
	}

	return nil
}
