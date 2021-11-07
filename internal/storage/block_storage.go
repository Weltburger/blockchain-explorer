package storage

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"explorer/models"
	"fmt"
	"log"
	"time"
)

type BlockStorage struct {
	database *Database
}

func (blockStorage *BlockStorage) Prepare() *sql.Stmt {
	st, err := blockStorage.database.Tx.Prepare(`
		INSERT INTO blocks.block (
			Protocol,
		    ChainID,
		    Hash,
		    Timestamp,
		    Header,
		    Metadata,
		    Operations
		) VALUES (
			?, ?, ?, ?, ?, ?, ?
		)`)

	if err != nil {
		log.Fatal(err)
	}

	return st
}

func (blockStorage *BlockStorage) Exc(data *models.Block) error {
	layout := "2006-01-02T15:04:05Z"
	timeStr := data.Header.Timestamp
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return err
	}

	/*ops, err := json.Marshal(data.Operations)
	if err != nil {
		return err
	}*/
	hdr, err := json.Marshal(data.Header)
	if err != nil {
		return err
	}
	mtdt, err := json.Marshal(data.Metadata)
	if err != nil {
		return err
	}

	if _, err := blockStorage.database.Stmt.Exec(
		data.Protocol,
		data.ChainID,
		data.Hash,
		t,
		string(hdr),
		string(mtdt),
		//string(ops),
	); err != nil {
		return err
	}

	return nil
}

func (blockStorage *BlockStorage) Cmt() error {
	if err := blockStorage.database.Tx.Commit(); err != nil {
		return err
	}

	blockStorage.database.Tx, _ = blockStorage.database.DB.Begin()
	blockStorage.database.Stmt = blockStorage.Prepare()

	return nil
}

func (blockStorage *BlockStorage) GetBlock(blk string) (*models.Block, error) {
	resp, err := blockStorage.database.DB.Query(`
		SELECT * FROM block WHERE Hash = ?
	`, blk)
	if err != nil {
		return nil, err
	}

	var (
		tm time.Time
		protocol, chainID, hash, header, metadata, ops string
		full []byte
	)

	resp.Next()
	err = resp.Scan(&protocol, &chainID, &hash, &tm, &header, &metadata, &ops)
	if err != nil {
		return nil, err
	}

	bts := [][]byte{[]byte("{\"protocol\":\"" + protocol + "\""),
		[]byte("\"chain_id\":\"" + chainID + "\""),
		[]byte("\"hash\":\"" + hash + "\""),
		[]byte("\"header\":" + header),
		[]byte("\"metadata\":" + metadata),
		[]byte("\"operations\":" + ops + "}")}

	full = bytes.Join(bts, []byte(","))

	block, err := models.UnmarshalBlock(full)
	if err != nil {
		return nil, err
	}

	return &block, nil
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

func (blockStorage *BlockStorage) GetBlocks(offset, limit int) ([]*models.Block, error) {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 1
	}

	resp, err := blockStorage.database.DB.Query(`
		SELECT * FROM block 
		LIMIT ?, ?
	`, offset, limit)
	if err != nil {
		return nil, err
	}

	var (
		tm time.Time
		protocol, chainID, hash, header, metadata, ops string
		full []byte
		blocks []*models.Block
	)

	for resp.Next() {
		err = resp.Scan(&protocol, &chainID, &hash, &tm, &header, &metadata, &ops)
		if err != nil {
			return nil, err
		}

		bts := [][]byte{[]byte("{\"protocol\":\"" + protocol + "\""),
			[]byte("\"chain_id\":\"" + chainID + "\""),
			[]byte("\"hash\":\"" + hash + "\""),
			[]byte("\"header\":" + header),
			[]byte("\"metadata\":" + metadata),
			[]byte("\"operations\":" + ops + "}")}

		full = bytes.Join(bts, []byte(","))

		block, err := models.UnmarshalBlock(full)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, &block)
	}

	fmt.Println(len(blocks))
	return blocks, nil
}
