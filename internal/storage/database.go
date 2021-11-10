package storage

import (
	"database/sql"
	_ "github.com/mailru/go-clickhouse"
	"log"
)

type Database struct {
	DB                 *sql.DB
	BlockTx            *sql.Tx
	BlockStmt          *sql.Stmt
	TransactionTx      *sql.Tx
	TransactionStmt    *sql.Stmt
	blockStorage       *BlockStorage
	transactionStorage *TransactionStorage
}

func (database *Database) BlockStorage() *BlockStorage {
	if database.blockStorage != nil {
		return database.blockStorage
	}

	database.blockStorage = &BlockStorage{database: database}

	return database.blockStorage
}

func (database *Database) TransactionStorage() *TransactionStorage {
	if database.transactionStorage != nil {
		return database.transactionStorage
	}

	database.transactionStorage = &TransactionStorage{database: database}

	return database.transactionStorage
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
