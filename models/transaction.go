package models

type Transaction struct {
	BlockHash        string `json:"block_hash"`
	Hash             string `json:"hash"`
	Branch           string `json:"branch"`
	Source           string `json:"source"`
	Destination      string `json:"destination"`
	Fee              string `json:"fee"`
	Counter          string `json:"counter"`
	GasLimit         string `json:"gas_limit"`
	Amount           string `json:"amount"`
	ConsumedMilligas string `json:"consumed_milligas"`
	StorageSize      string `json:"storage_size"`
	Signature        string `json:"signature"`
}

type TransactionMainInfo struct {
	BlockHash        string `json:"block_hash"`
	Hash             string `json:"hash"`
	Source           string `json:"source"`
	Destination      string `json:"destination"`
	Fee              string `json:"fee"`
	Amount           string `json:"amount"`
}
