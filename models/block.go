package models

import (
	"encoding/json"
)

func UnmarshalBlock(data []byte) (Block, error) {
	var r Block
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Block) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Block struct {
	Protocol   Protocol      `json:"protocol"`
	ChainID    ChainID       `json:"chain_id"`
	Hash       string        `json:"hash"`
	Header     Header        `json:"header"`
	Metadata   BlockMetadata `json:"metadata"`
	Operations [][]Operation `json:"operations"`
}
