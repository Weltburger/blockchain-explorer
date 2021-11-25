package usecase

import (
	"context"
	"explorer/internal/explorer"
	"explorer/models"
)

type TransUseCase struct {
	transRepo explorer.TransRepo
}

func NewTransUseCase(transRepo explorer.TransRepo) *TransUseCase {
	return &TransUseCase{transRepo: transRepo}
}

func (t *TransUseCase) GetTransactions(ctx context.Context, offset, limit int, blk, hash, acc string) ([]models.Transaction, error) {
	transactions, err := t.transRepo.GetTransactions(ctx, offset, limit, blk, hash, acc)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
