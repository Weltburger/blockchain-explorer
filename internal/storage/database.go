package storage

import (
	"database/sql"
	_ "github.com/mailru/go-clickhouse"
	"log"
)

type Database struct {
	DB *sql.DB
}

func GetDB() *Database {
	DB := startDB()

	return &Database{DB: DB}
}

func startDB() *sql.DB {
	db, err := sql.Open("clickhouse", "http://localhost:8123/blocks?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
