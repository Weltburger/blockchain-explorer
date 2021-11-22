package storage

import (
	"database/sql"
	"explorer/models"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type TransactionStorage struct {
	database *Database
	Tx       *sql.Tx
	Stmt     *sql.Stmt
}

func (transactionStorage *TransactionStorage) PrepareTransactionTx() error {
	trx, err := transactionStorage.database.DB.Begin()
	if err != nil {
		return err
	}

	st, err := trx.Prepare(`
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
		return err
	}

	transactionStorage.Stmt, transactionStorage.Tx = st, trx

	return nil
}

func (transactionStorage *TransactionStorage) Exc(data *models.Transaction) error {
	if _, err := transactionStorage.Stmt.Exec(
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
		_ = transactionStorage.Stmt.Close()
		_ = transactionStorage.Tx.Rollback()
		return err
	}

	return nil
}

func (transactionStorage *TransactionStorage) Cmt() error {
	if err := transactionStorage.Tx.Commit(); err != nil {
		_ = transactionStorage.Tx.Rollback()
		return err
	}

	return nil
}

func (transactionStorage *TransactionStorage) SaveTransaction(transaction *models.Transaction) error {
	err := transactionStorage.Exc(transaction)
	if err != nil {
		return err
	}

	return nil
}

func (transactionStorage *TransactionStorage) GetTransactions(offset, limit int, blk, hash, acc string) ([]models.Transaction, error) {
	if blk != "" {
		blk = fmt.Sprintf("block_hash='%s'", blk)
	}
	if hash != "" {
		hash = fmt.Sprintf("hash='%s'", hash)
	}
	if acc != "" {
		acc = fmt.Sprintf("(source='%s' OR destination='%s')", acc, acc)
	}
	query, _, err := sq.Select("block_hash",
								 "hash",
								 "branch",
								 "source",
								 "destination",
								 "fee, counter",
								 "gas_limit",
								 "amount",
								 "consumed_milligas",
								 "storage_size",
								 "signature").
		From("transactions").
		Where(blk).
		Where(hash).
		Where(acc).
		Limit(uint64(limit)).
		Offset(uint64(offset)).ToSql()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := transactionStorage.database.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	transaction := new(models.Transaction)

	for resp.Next() {
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

		transactions = append(transactions, *transaction)
	}

	return transactions, nil
}
