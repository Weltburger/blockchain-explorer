package explorer

import (
	"context"
	"explorer/models"
)

type BlockRepo interface {
	Cmt() error
	PrepareBlockTx() error
	Exc(data *models.Block) error
	SaveBlock(block *models.Block) error
	GetBlock(ctx context.Context, blk string) (*models.Block, error)
	GetBlocks(ctx context.Context, offset, limit int) ([]models.Block, error)
}

type TransRepo interface {
	Cmt() error
	PrepareTransactionTx() error
	Exc(data *models.Transaction) error
	SaveTransaction(transaction *models.Transaction) error
	GetTransactions(ctx context.Context,offset, limit int, blk, hash, acc string) ([]models.Transaction, error)
}
