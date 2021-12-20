package explorer

import (
	"context"
	"explorer/models"
)

type BlockUseCase interface {
	GetBlock(ctx context.Context, blk string) (*models.Block, error)
	GetBlocks(ctx context.Context, offset, limit int) ([]models.Block, error)
}

type TransUseCase interface {
	GetTransactions(ctx context.Context, offset, limit int, blk, hash, acc string) ([]models.Transaction, error)
}

type TransMIUseCase interface {
	GetTransactions(ctx context.Context, offset, limit int, blk, hash, acc string) ([]models.TransactionMainInfo, error)
}