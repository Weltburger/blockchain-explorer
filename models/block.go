package models

import "encoding/json"

func UnmarshalBlock(data []byte) (Block, error) {
	var r Block
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Block) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Block struct {
	Protocol   string        `json:"protocol"`
	ChainID    string        `json:"chain_id"`
	Hash       string        `json:"hash"`
	Header     Header        `json:"header"`
	Metadata   BlockMetadata `json:"metadata"`
	Operations []byte		 `json:"operations"`
}

type Header struct {
	Proto                     int64       `json:"proto"`
	Predecessor               string      `json:"predecessor"`
	Timestamp                 string      `json:"timestamp"`
	ValidationPass            int64       `json:"validation_pass"`
	OperationsHash            string      `json:"operations_hash"`
	Fitness                   []string    `json:"fitness"`
	Priority                  int64       `json:"priority"`
	Signature                 string      `json:"signature"`
	BakerFee                  uint64      `json:"baker_fee"`
}

type BlockMetadata struct {
	Baker                     string                     `json:"baker"`
	LevelInfo                 LevelInfo                  `json:"level_info"`
	ConsumedGas               string                     `json:"consumed_gas"`
}

type LevelInfo struct {
	Level              int64 `json:"level"`
	Cycle              int64 `json:"cycle"`
	CyclePosition      int64 `json:"cycle_position"`
}
