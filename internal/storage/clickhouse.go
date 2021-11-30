package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mailru/go-clickhouse"
)

type ClickhouseDataSource struct {
	DB *sql.DB
}

func InitClickhouse() (*ClickhouseDataSource, error) {
	chHost := os.Getenv("CH_HOST")
	chPort := os.Getenv("CH_PORT")
	chDB := os.Getenv("CH_DB")
	chDebug := os.Getenv("CH_DEBUG")

	chConnString := fmt.Sprintf("http://%s:%s/%s?debug=%s", chHost, chPort, chDB, chDebug)

	db, err := sql.Open("clickhouse", chConnString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &ClickhouseDataSource{DB: db}, nil
}
