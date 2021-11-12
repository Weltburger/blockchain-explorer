package storage

import (
	"database/sql"
	"explorer/models"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type TransactionStorage struct {
	database *Database
	Tx       *sql.Tx
	Stmt     *sql.Stmt
}

func PrepareTransaction(tx *sql.Tx) *sql.Stmt {
	st, err := tx.Prepare(`
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
		return err
	}

	return nil
}

func (transactionStorage *TransactionStorage) Cmt() error {
	if err := transactionStorage.Tx.Commit(); err != nil {
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

func (transactionStorage *TransactionStorage) GetTransactions(offset, limit int, blk, hash, acc string) ([]*models.Transaction, error) {
	if blk != "" {
		blk = "block_hash='"+ blk +"'"
	}
	if hash != "" {
		hash = "hash='"+hash+"'"
	}
	if acc != "" {
		acc = "(source='" + acc + "' OR destination='"+acc + "')"
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
