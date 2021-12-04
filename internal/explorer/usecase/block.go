package usecase

import (
	"context"
	"explorer/internal/explorer"
	"explorer/models"
	"strconv"
)

type BlockUseCase struct {
	blockRepo explorer.BlockRepo
}

func NewBlockUseCase(blockRepo explorer.BlockRepo) *BlockUseCase {
	return &BlockUseCase{blockRepo: blockRepo}
}

func (b *BlockUseCase) GetBlock(ctx context.Context, blk string) (*models.Block, error) {
	blkLevel, err := strconv.ParseInt(blk, 10, 64)
	if err != nil {
		block, err := b.blockRepo.GetBlockByHash(ctx, blk)
		if err != nil {
			return nil, err
		}
		return block, nil
	}

	block, err := b.blockRepo.GetBlockByLevel(ctx, blkLevel)
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
