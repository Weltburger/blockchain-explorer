package explorer

import (
	"context"
	"explorer/models"
)

type BlockRepo interface {
	GetBlockByHash(ctx context.Context, blk string) (*models.Block, error)
	GetBlockByLevel(ctx context.Context, blk int64) (*models.Block, error)
	GetBlocks(ctx context.Context, offset, limit int) ([]models.Block, error)
}

type TransRepo interface {
	GetTransactions(ctx context.Context,offset, limit int, blk, hash, acc string) ([]models.Transaction, error)
}
