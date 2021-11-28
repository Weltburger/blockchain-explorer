package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/mailru/go-clickhouse"
	"github.com/spf13/viper"
)

type ClickhouseDataSource struct {
	DB *sql.DB
}

func InitClickhouse() (*ClickhouseDataSource, error) {
	chHost := viper.GetString("clickhouse.ch_host")
	chPort := viper.GetInt("clickhouse.ch_port")
	chDB := viper.GetString("clickhouse.ch_db")
	chDebug := viper.GetString("clickhouse.ch_debug")

	chConnString := fmt.Sprintf("http://%s:%d/%s?debug=%s", chHost, chPort, chDB, chDebug)

	db, err := sql.Open("clickhouse", chConnString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &ClickhouseDataSource{DB: db}, nil
}
