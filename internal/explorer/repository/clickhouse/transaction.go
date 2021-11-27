package clickhouse

import (
	"context"
	"database/sql"
	"explorer/internal/apperrors"
	"explorer/models"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type TransactionRepository struct {
	DB       *sql.DB
	Tx       *sql.Tx
	Stmt     *sql.Stmt
}

func NewTransRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		DB: db,
	}
}

func (t *TransactionRepository) PrepareTransactionTx() error {
	trx, err := t.DB.Begin()
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
		trx.Rollback()
		return err
	}

	t.Stmt, t.Tx = st, trx

	return nil
}

func (t *TransactionRepository) Exc(data *models.Transaction) error {
	if _, err := t.Stmt.Exec(
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

func (t *TransactionRepository) Cmt() error {
	if err := t.Tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepository) SaveTransaction(transaction *models.Transaction) error {
	err := t.Exc(transaction)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionRepository) GetTransactions(ctx context.Context,offset, limit int, blk, hash, acc string) ([]models.Transaction, error) {
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
		return nil, err
	}

	resp, err := t.DB.QueryContext(ctx, query)
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
			return nil, apperrors.NewNotFound("clickhouse", "such transactions was")
		}

		transactions = append(transactions, *transaction)
	}

	return transactions, nil
}
