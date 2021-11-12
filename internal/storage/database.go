package storage

import (
	"database/sql"
	_ "github.com/mailru/go-clickhouse"
	"log"
)

type Database struct {
	DB *sql.DB
}

func (database *Database) BlockStorage() *BlockStorage {
	trx, _ := database.DB.Begin()
	return &BlockStorage{
		database: database,
		Tx:       trx,
		Stmt:     PrepareBlock(trx),
	}
}

func (database *Database) TransactionStorage() *TransactionStorage {
	trx, _ := database.DB.Begin()
	return &TransactionStorage{
		database: database,
		Tx:       trx,
		Stmt:     PrepareTransaction(trx),
	}
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

func (database *Database) CloseDB() {
	database.DB.Close()
}
