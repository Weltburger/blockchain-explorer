package clickhouse

import (
	"context"
	"database/sql"
	"explorer/internal/apperrors"
	"explorer/models"
	"strings"
	"time"
)

type BlockRepository struct {
	DB       *sql.DB
	Tx       *sql.Tx
	Stmt     *sql.Stmt
}

func NewBlockRepository(db *sql.DB) *BlockRepository {
	return &BlockRepository{
		DB: db,
	}
}

func (b *BlockRepository) PrepareBlockTx() error {
	trx, err := b.DB.Begin()
	if err != nil {
		return err
	}

	st, err := trx.Prepare(`
		INSERT INTO blocks.block (
			protocol,
			chain_id,
			hash,
			baker_fees,
			"level",
			predecessor,
			priority,
			"timestamp",
			validation_pass,
			validation_hash,
			fitness,
			signature,
			baker,
			cycle_num,
			cycle_position,
			consumed_gas     
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`)
	if err != nil {
		trx.Rollback()
		return err
	}

	b.Stmt, b.Tx = st, trx

	return nil
}

func (b *BlockRepository) Exc(data *models.Block) error {
	layout := "2006-01-02T15:04:05Z"
	timeStr := data.Header.Timestamp
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return err
	}

	fitness := strings.Join(data.Header.Fitness, ",")

	if _, err := b.Stmt.Exec(
		data.Protocol,
		data.ChainID,
		data.Hash,
		0,
		data.Metadata.LevelInfo.Level,
		data.Header.Predecessor,
		data.Header.Priority,
		t,
		data.Header.ValidationPass,
		data.Header.OperationsHash,
		fitness,
		data.Header.Signature,
		data.Metadata.Baker,
		data.Metadata.LevelInfo.Cycle,
		data.Metadata.LevelInfo.CyclePosition,
		data.Metadata.ConsumedGas,
	); err != nil {
		return err
	}

	return nil
}

func (b *BlockRepository) Cmt() error {
	if err := b.Tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (b *BlockRepository) Rollback() error {
	if err := b.Tx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (b *BlockRepository) GetBlockByHash(ctx context.Context, blk string) (*models.Block, error) {
	resp, err := b.DB.QueryContext(ctx, `
		SELECT protocol,
			   chain_id,
			   hash,
			   baker_fees,
			   "level",
			   predecessor,
			   priority,
			   "timestamp",
			   validation_pass,
			   validation_hash,
			   fitness,
			   signature,
			   baker,
			   cycle_num,
			   cycle_position,
			   consumed_gas 
		FROM blocks.block WHERE hash = ?
	`, blk)
	if err != nil {
		return nil, err
	}

	var tm time.Time
	var fitness string
	block := new(models.Block)

	resp.Next()
	err = resp.Scan(&block.Protocol,
		&block.ChainID,
		&block.Hash,
		&block.Header.BakerFee,
		&block.Metadata.LevelInfo.Level,
		&block.Header.Predecessor,
		&block.Header.Priority,
		&tm,
		&block.Header.ValidationPass,
		&block.Header.OperationsHash,
		&fitness,
		&block.Header.Signature,
		&block.Metadata.Baker,
		&block.Metadata.LevelInfo.Cycle,
		&block.Metadata.LevelInfo.CyclePosition,
		&block.Metadata.ConsumedGas)
	if err != nil {
		return nil, apperrors.NewNotFound("clickhouse", "such block was")
	}

	block.Header.Timestamp = tm.String()
	block.Header.Fitness = strings.Split(fitness, ",")

	return block, nil
}

func (b *BlockRepository) GetBlockByLevel(ctx context.Context, blk int64) (*models.Block, error) {
	resp, err := b.DB.QueryContext(ctx, `
		SELECT protocol,
			   chain_id,
			   hash,
			   baker_fees,
			   "level",
			   predecessor,
			   priority,
			   "timestamp",
			   validation_pass,
			   validation_hash,
			   fitness,
			   signature,
			   baker,
			   cycle_num,
			   cycle_position,
			   consumed_gas 
		FROM blocks.block WHERE "level" = ?
	`, blk)
	if err != nil {
		return nil, err
	}

	var tm time.Time
	var fitness string
	block := new(models.Block)

	resp.Next()
	err = resp.Scan(&block.Protocol,
		&block.ChainID,
		&block.Hash,
		&block.Header.BakerFee,
		&block.Metadata.LevelInfo.Level,
		&block.Header.Predecessor,
		&block.Header.Priority,
		&tm,
		&block.Header.ValidationPass,
		&block.Header.OperationsHash,
		&fitness,
		&block.Header.Signature,
		&block.Metadata.Baker,
		&block.Metadata.LevelInfo.Cycle,
		&block.Metadata.LevelInfo.CyclePosition,
		&block.Metadata.ConsumedGas)
	if err != nil {
		return nil, apperrors.NewNotFound("clickhouse", "such block was")
	}

	block.Header.Timestamp = tm.String()
	block.Header.Fitness = strings.Split(fitness, ",")

	return block, nil
}

func (b *BlockRepository) SaveBlock(block *models.Block) error {
	err := b.Exc(block)
	if err != nil {
		return err
	}
	err = b.Cmt()
	if err != nil {
		return err
	}

	return nil
}

func (b *BlockRepository) GetBlocks(ctx context.Context, offset, limit int) ([]models.Block, error) {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 1
	}

	resp, err := b.DB.QueryContext(ctx,`
		SELECT protocol,
			   chain_id,
			   hash,
			   baker_fees,
			   "level",
			   predecessor,
			   priority,
			   "timestamp",
			   validation_pass,
			   validation_hash,
			   fitness,
			   signature,
			   baker,
			   cycle_num,
			   cycle_position,
			   consumed_gas 
		FROM blocks.block 
		LIMIT ?, ?
	`, offset, limit)
	if err != nil {
		return nil, apperrors.NewNotFound("clickhouse", "such blocks was")
	}

	var (
		tm      time.Time
		fitness string
		blocks  []models.Block
	)

	block := new(models.Block)
	for resp.Next() {
		err = resp.Scan(&block.Protocol,
			&block.ChainID,
			&block.Hash,
			&block.Header.BakerFee,
			&block.Metadata.LevelInfo.Level,
			&block.Header.Predecessor,
			&block.Header.Priority,
			&tm,
			&block.Header.ValidationPass,
			&block.Header.OperationsHash,
			&fitness,
			&block.Header.Signature,
			&block.Metadata.Baker,
			&block.Metadata.LevelInfo.Cycle,
			&block.Metadata.LevelInfo.CyclePosition,
			&block.Metadata.ConsumedGas)
		if err != nil {
			return nil, apperrors.NewNotFound("clickhouse", "such blocks was")
		}

		block.Header.Timestamp = tm.String()
		block.Header.Fitness = strings.Split(fitness, ",")

		blocks = append(blocks, *block)
	}

	return blocks, nil
}
