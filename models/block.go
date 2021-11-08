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
	Operations [][]Operation `json:"operations"`
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

//////////////////////////////////////////////////////////


type Operation struct {
	Protocol  string    `json:"protocol"`
	ChainID   string    `json:"chain_id"`
	Hash      string    `json:"hash"`
	Branch    string    `json:"branch"`
	Contents  []Content `json:"contents"`
	Signature string    `json:"signature"`
}

type Content struct {
	Kind         string     `json:"kind"`
	Source       string     `json:"source"`
	Fee          string     `json:"fee"`
	Counter      string     `json:"counter"`
	GasLimit     string     `json:"gas_limit"`
	StorageLimit string     `json:"storage_limit"`
	Amount       string     `json:"amount"`
	Destination  string     `json:"destination"`
	Parameters   Parameters `json:"parameters"`
	Metadata     Metadata   `json:"metadata"`
}

type Metadata struct {
	BalanceUpdates           []MetadataBalanceUpdate   `json:"balance_updates"`
	OperationResult          OperationResult           `json:"operation_result"`
	InternalOperationResults []InternalOperationResult `json:"internal_operation_results"`
}

type MetadataBalanceUpdate struct {
	Kind     string  `json:"kind"`
	Contract *string `json:"contract,omitempty"`
	Change   string  `json:"change"`
	Origin   string  `json:"origin"`
	Category *string `json:"category,omitempty"`
	Delegate *string `json:"delegate,omitempty"`
	Cycle    *int64  `json:"cycle,omitempty"`
}

type InternalOperationResult struct {
	Kind        string `json:"kind"`
	Source      string `json:"source"`
	Nonce       int64  `json:"nonce"`
	Amount      string `json:"amount"`
	Destination string `json:"destination"`
	Result      Result `json:"result"`
}

type Result struct {
	Status           string                `json:"status"`
	BalanceUpdates   []ResultBalanceUpdate `json:"balance_updates"`
	ConsumedGas      string                `json:"consumed_gas"`
	ConsumedMilligas string                `json:"consumed_milligas"`
}

type ResultBalanceUpdate struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract"`
	Change   string `json:"change"`
	Origin   string `json:"origin"`
}

type OperationResult struct {
	Status           string            `json:"status"`
	Storage          []Storage         `json:"storage"`
	BigMapDiff       []BigMapDiff      `json:"big_map_diff"`
	ConsumedGas      string            `json:"consumed_gas"`
	ConsumedMilligas string            `json:"consumed_milligas"`
	StorageSize      string            `json:"storage_size"`
	LazyStorageDiff  []LazyStorageDiff `json:"lazy_storage_diff"`
}

type BigMapDiff struct {
	Action  string `json:"action"`
	BigMap  string `json:"big_map"`
	KeyHash string `json:"key_hash"`
	Key     Value  `json:"key"`
}

type Value struct {
	Int string `json:"int"`
}

type LazyStorageDiff struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
	Diff Diff   `json:"diff"`
}

type Diff struct {
	Action  string   `json:"action"`
	Updates []Update `json:"updates"`
}

type Update struct {
	KeyHash string `json:"key_hash"`
	Key     Value  `json:"key"`
}

type Storage struct {
	Prim *string `json:"prim,omitempty"`
	Args []Arg   `json:"args,omitempty"`
	Int  *string `json:"int,omitempty"`
}

type Arg struct {
	Bytes *string `json:"bytes,omitempty"`
	Prim  *string `json:"prim,omitempty"`
	Args  []Value `json:"args,omitempty"`
	Int   *string `json:"int,omitempty"`
}

type Parameters struct {
	Entrypoint string `json:"entrypoint"`
	Value      Value  `json:"value"`
}

