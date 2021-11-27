package usecase

import (
	"context"
	"explorer/internal/explorer"
	"explorer/models"
)

type BlockUseCase struct {
	blockRepo explorer.BlockRepo
}

func NewBlockUseCase(blockRepo explorer.BlockRepo) *BlockUseCase {
	return &BlockUseCase{blockRepo: blockRepo}
}

func (b *BlockUseCase) GetBlock(ctx context.Context, blk string) (*models.Block, error) {
	block, err := b.blockRepo.GetBlock(ctx, blk)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (b *BlockUseCase) GetBlocks(ctx context.Context, offset, limit int) ([]models.Block, error) {
	blocks, err := b.blockRepo.GetBlocks(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return blocks, nil
}
