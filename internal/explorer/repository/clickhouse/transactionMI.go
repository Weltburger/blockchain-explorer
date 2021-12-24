package clickhouse

import (
	"context"
	"database/sql"
	"explorer/internal/apperrors"
	"explorer/models"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type TransactionMIRepository struct {
	DB       *sql.DB
	Tx       *sql.Tx
	Stmt     *sql.Stmt
}

func NewTransMIRepository(db *sql.DB) *TransactionMIRepository {
	return &TransactionMIRepository{
		DB: db,
	}
}

func (t *TransactionMIRepository) PrepareTransactionTx() error {
	trx, err := t.DB.Begin()
	if err != nil {
		return err
	}

	st, err := trx.Prepare(`
		INSERT INTO blocks.transactions_mi (
			block_hash, 
			hash, 
			source, 
			destination, 
			fee, 
			amount
		) VALUES (
			?, ?, ?, ?, ?, ?
		)`)
	if err != nil {
		trx.Rollback()
		return err
	}

	t.Stmt, t.Tx = st, trx

	return nil
}

func (t *TransactionMIRepository) Exc(data *models.TransactionMainInfo) error {
	if _, err := t.Stmt.Exec(
		data.BlockHash,
		data.Hash,
		data.Source,
		data.Destination,
		data.Fee,
		data.Amount,
	); err != nil {
		return err
	}

	return nil
}

func (t *TransactionMIRepository) Cmt() error {
	if err := t.Tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (t *TransactionMIRepository) SaveTransaction(transaction *models.TransactionMainInfo) error {
	err := t.Exc(transaction)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionMIRepository) GetTransactions(ctx context.Context, offset, limit int, blk, hash, acc string) ([]models.TransactionMainInfo, error) {
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
		"source",
		"destination",
		"fee",
		"amount").
		From("transactions_mi").
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

	var transactions []models.TransactionMainInfo
	transaction := new(models.TransactionMainInfo)

	for resp.Next() {
		err = resp.Scan(&transaction.BlockHash,
			&transaction.Hash,
			&transaction.Source,
			&transaction.Destination,
			&transaction.Fee,
			&transaction.Amount)
		if err != nil {
			return nil, apperrors.NewNotFound("clickhouse", "such transactionsMI was")
		}

		transactions = append(transactions, *transaction)
	}

	return transactions, nil
}
