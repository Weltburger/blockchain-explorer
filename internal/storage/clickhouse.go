package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mailru/go-clickhouse"
)

type ClickhouseDataSource struct {
	DB *sql.DB
}

func InitClickhouse() (*ClickhouseDataSource, error) {
	chHost := os.Getenv("CLICKHOUSE_HOST")
	chPort := os.Getenv("CLICKHOUSE_PORT")
	chDB := os.Getenv("CLICKHOUSE_DB")
	chDebug := os.Getenv("CLICKHOUSE_DEBUG")

	chConnString := fmt.Sprintf("http://%s:%s/%s?debug=%s", chHost, chPort, chDB, chDebug)

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

	sourceURL, err := (&file.File{}).Open("file://internal/storage/migrations/ch")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithInstance(
		"file", sourceURL,
		"clickhouse", driver)
	if err != nil {
		return nil, fmt.Errorf("error create migrating: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	log.Printf("Clickhouse Migration done\n")

	return &ClickhouseDataSource{DB: db}, nil
}
