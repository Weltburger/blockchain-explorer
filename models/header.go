package models

type Header struct {
	Level                     uint64      `json:"level"`
	Proto                     int64       `json:"proto"`
	Predecessor               Predecessor `json:"predecessor"`
	Timestamp                 string      `json:"timestamp"`
	ValidationPass            int64       `json:"validation_pass"`
	OperationsHash            string      `json:"operations_hash"`
	Fitness                   []string    `json:"fitness"`
	Context                   string      `json:"context"`
	Priority                  int64       `json:"priority"`
	ProofOfWorkNonce          string      `json:"proof_of_work_nonce"`
	LiquidityBakingEscapeVote bool        `json:"liquidity_baking_escape_vote"`
	Signature                 string      `json:"signature"`
}
