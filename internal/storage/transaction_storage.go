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

func (transactionStorage *TransactionStorage) GetTransactionsByBlock(block string) ([]*models.Transaction, error) {
	resp, err := transactionStorage.database.DB.Query(`
		SELECT * FROM transactions 
		WHERE block_hash = ?
	`, block)
	if err != nil {
		return nil, err
	}

	var transactions []*models.Transaction

	for resp.Next() {
		transaction := new(models.Transaction)
		err = resp.Scan(&transaction.BlockHash,
			&transaction.Hash,
			&transaction.Branch,
			&transaction.Source,
			&transaction.Destination,
			&transaction.Fee,
			&transaction.Counter,
			&transaction.GasLimit,
			&transaction.Amount,
			&transaction.ConsumedMilligas,
			&transaction.StorageSize,
			&transaction.Signature)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (transactionStorage *TransactionStorage) GetTransactions(offset, limit int) ([]*models.Transaction, error) {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 1
	}

	resp, err := transactionStorage.database.DB.Query(`
		SELECT * FROM transactions 
		LIMIT ?, ?
	`, offset, limit)
	if err != nil {
		return nil, err
	}

	var transactions []*models.Transaction

	for resp.Next() {
		transaction := new(models.Transaction)
		err = resp.Scan(&transaction.BlockHash,
			&transaction.Hash,
			&transaction.Branch,
			&transaction.Source,
			&transaction.Destination,
			&transaction.Fee,
			&transaction.Counter,
			&transaction.GasLimit,
			&transaction.Amount,
			&transaction.ConsumedMilligas,
			&transaction.StorageSize,
			&transaction.Signature)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (transactionStorage *TransactionStorage) GetTransactionsByAddress(address string) ([]*models.Transaction, error) {
	resp, err := transactionStorage.database.DB.Query(`
		SELECT * FROM transactions 
		WHERE source = ? OR destination = ?
	`, address, address)
	if err != nil {
		return nil, err
	}

	var transactions []*models.Transaction

	for resp.Next() {
		transaction := new(models.Transaction)
		err = resp.Scan(&transaction.BlockHash,
			&transaction.Hash,
			&transaction.Branch,
			&transaction.Source,
			&transaction.Destination,
			&transaction.Fee,
			&transaction.Counter,
			&transaction.GasLimit,
			&transaction.Amount,
			&transaction.ConsumedMilligas,
			&transaction.StorageSize,
			&transaction.Signature)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (transactionStorage *TransactionStorage) GetTransactionsByHash(hash string) ([]*models.Transaction, error) {
	resp, err := transactionStorage.database.DB.Query(`
		SELECT * FROM transactions 
		WHERE hash = ? 
	`, hash)
	if err != nil {
		return nil, err
	}

	var transactions []*models.Transaction

	for resp.Next() {
		transaction := new(models.Transaction)
		err = resp.Scan(&transaction.BlockHash,
			&transaction.Hash,
			&transaction.Branch,
			&transaction.Source,
			&transaction.Destination,
			&transaction.Fee,
			&transaction.Counter,
			&transaction.GasLimit,
			&transaction.Amount,
			&transaction.ConsumedMilligas,
			&transaction.StorageSize,
			&transaction.Signature)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
