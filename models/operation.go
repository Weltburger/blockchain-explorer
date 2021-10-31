package models

type Operation struct {
	Protocol  Protocol    `json:"protocol"`
	ChainID   ChainID     `json:"chain_id"`
	Hash      string      `json:"hash"`
	Branch    Predecessor `json:"branch"`
	Contents  []Content   `json:"contents"`
	Signature string      `json:"signature,omitempty"`
}
