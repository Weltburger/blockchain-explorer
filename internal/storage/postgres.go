package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDataSource struct {
	DB *sqlx.DB
}

// InitDS establishes connections to fields in PostgresDataSource
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
