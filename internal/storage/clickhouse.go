package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/mailru/go-clickhouse"
)

type ClickhouseDataSource struct {
	DB *sql.DB
}

func InitClickhouse() (*ClickhouseDataSource, error) {
	chHost := os.Getenv("CLICKHOUSE_HOST")
	chPort := os.Getenv("CLICKHOUSE_PORT")
	// chUser := os.Getenv("CLICKHOUSE_USER")
	// chPass := os.Getenv("CLICKHOUSE_PASSWORD")
	chDB := os.Getenv("CLICKHOUSE_DB")
	chDebug := os.Getenv("CLICKHOUSE_DEBUG")

	chConnString := fmt.Sprintf("http://%s:%s/%s?debug=%s", chHost, chPort, chDB, chDebug)

	log.Println(chConnString)
	db, err := sql.Open("clickhouse", chConnString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	driver, err := clickhouse.WithInstance(db, &clickhouse.Config{})
	if err != nil {
		return nil, fmt.Errorf("error create driver for migrate: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../explorer/migration",
		"clickhouse", driver)
	if err != nil {
		return nil, fmt.Errorf("error create migrating: %v", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	log.Printf("Migrate test data\n")

	return &ClickhouseDataSource{DB: db}, nil
}
