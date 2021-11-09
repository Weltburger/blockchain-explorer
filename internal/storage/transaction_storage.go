package storage

import (
	"database/sql"
	"explorer/models"
	"log"
)

type TransactionStorage struct {
	database *Database
}

func (transactionStorage *TransactionStorage) Prepare() *sql.Stmt {
	st, err := transactionStorage.database.TransactionTx.Prepare(`
		INSERT INTO blocks.transactions (
			block_hash, 
			hash, 
			branch, 
			source, 
			destination, 
			fee, 
			counter, 
			gas_limit, 
			amount, 
			consumed_milligas, 
			storage_size, 
			signature
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`)

	if err != nil {
		log.Fatal(err)
	}

	return st
}

func (transactionStorage *TransactionStorage) Exc(data *models.Transaction) error {
	if _, err := transactionStorage.database.TransactionStmt.Exec(
		data.BlockHash,
		data.Hash,
		data.Branch,
		data.Source,
		data.Destination,
		data.Fee,
		data.Counter,
		data.GasLimit,
		data.Amount,
		data.ConsumedMilligas,
		data.StorageSize,
		data.Signature,
	); err != nil {
		return err
	}

	return nil
}

func (transactionStorage *TransactionStorage) Cmt() error {
	if err := transactionStorage.database.TransactionTx.Commit(); err != nil {
		return err
	}

	transactionStorage.database.TransactionTx, _ = transactionStorage.database.DB.Begin()
	transactionStorage.database.TransactionStmt = transactionStorage.Prepare()

	return nil
}

func (transactionStorage *TransactionStorage) SaveTransaction(transaction *models.Transaction) error {
	err := transactionStorage.Exc(transaction)
	if err != nil {
		return err
	}

	return nil
}
