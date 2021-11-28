package storage

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type PostgresDataSource struct {
	DB *sqlx.DB
}

// InitDS establishes connections to fields in PostgresDataSource
func InitPostgres() (*PostgresDataSource, error) {
	log.Printf("Initializing data sources\n")

	pgHost := viper.GetString("postgres.pg_host")
	pgPort := viper.GetInt("postgres.pg_port")
	pgUser := viper.GetString("postgres.pg_user")
	pgPassword := viper.GetString("postgres.pg_pass")
	pgDB := viper.GetString("postgres.pg_db")
	pgSSL := viper.GetString("postgres.pg_ssl")

	pgConnString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL)

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
