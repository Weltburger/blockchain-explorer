package storage

import (
	"database/sql"
	"encoding/json"
	"explorer/models"
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

func (blockStorage *BlockStorage) Exc(data *models.Block) {
	layout := "2006-01-02T15:04:05Z"
	timeStr := data.Header.Timestamp
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Fatal(err)
	}

	ops, err := json.Marshal(data.Operations)
	if err != nil {
		log.Fatal(err)
	}
	hdr, err := json.Marshal(data.Header)
	if err != nil {
		log.Fatal(err)
	}
	mtdt, err := json.Marshal(data.Metadata)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := blockStorage.database.Stmt.Exec(
		data.Protocol,
		data.ChainID,
		data.Hash,
		t,
		string(hdr),
		string(mtdt),
		string(ops),
	); err != nil {
		log.Fatal(err)
	}
}

func (blockStorage *BlockStorage) Cmt() {
	if err := blockStorage.database.Tx.Commit(); err != nil {
		log.Fatal(err)
	}

	blockStorage.database.Tx, _ = blockStorage.database.DB.Begin()
	blockStorage.database.Stmt = blockStorage.Prepare()
}

