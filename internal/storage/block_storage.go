package storage

import (
	"context"
	"database/sql"
	"explorer/models"
	"strings"
	"time"
)

type BlockStorage struct {
	database *Database
	Tx       *sql.Tx
	Stmt     *sql.Stmt
}

func (blockStorage *BlockStorage) PrepareBlockTx() error {
	trx, err := blockStorage.database.DB.Begin()
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

	blockStorage.Stmt, blockStorage.Tx = st, trx

	return nil
}

func (blockStorage *BlockStorage) Exc(data *models.Block) error {
	layout := "2006-01-02T15:04:05Z"
	timeStr := data.Header.Timestamp
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return err
	}

	fitness := strings.Join(data.Header.Fitness, ",")

	if _, err := blockStorage.Stmt.Exec(
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

func (blockStorage *BlockStorage) Cmt() error {
	if err := blockStorage.Tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (blockStorage *BlockStorage) GetBlock(ctx context.Context, blk string) (*models.Block, error) {
	resp, err := blockStorage.database.DB.QueryContext(ctx, `
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
		FROM block WHERE hash = ?
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
		return nil, err
	}

	block.Header.Timestamp = tm.String()
	block.Header.Fitness = strings.Split(fitness, ",")

	return block, nil
}

func (blockStorage *BlockStorage) SaveBlock(block *models.Block) error {
	err := blockStorage.Exc(block)
	if err != nil {
		return err
	}
	err = blockStorage.Cmt()
	if err != nil {
		return err
	}

	return nil
}

func (blockStorage *BlockStorage) GetBlocks(ctx context.Context, offset, limit int) ([]models.Block, error) {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 1
	}

	resp, err := blockStorage.database.DB.QueryContext(ctx,`
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
		FROM block 
		LIMIT ?, ?
	`, offset, limit)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		block.Header.Timestamp = tm.String()
		block.Header.Fitness = strings.Split(fitness, ",")

		blocks = append(blocks, *block)
	}

	return blocks, nil
}
