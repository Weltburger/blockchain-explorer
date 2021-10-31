package models

type BlockMetadata struct {
	Protocol                  Protocol                   `json:"protocol"`
	NextProtocol              Protocol                   `json:"next_protocol"`
	TestChainStatus           TestChainStatus            `json:"test_chain_status"`
	MaxOperationsTTL          int64                      `json:"max_operations_ttl"`
	MaxOperationDataLength    int64                      `json:"max_operation_data_length"`
	MaxBlockHeaderLength      int64                      `json:"max_block_header_length"`
	MaxOperationListLength    []MaxOperationListLength   `json:"max_operation_list_length"`
	Baker                     string                     `json:"baker"`
	LevelInfo                 LevelInfo                  `json:"level_info"`
	VotingPeriodInfo          VotingPeriodInfo           `json:"voting_period_info"`
	NonceHash                 interface{}                `json:"nonce_hash"`
	ConsumedGas               string                     `json:"consumed_gas"`
	Deactivated               []interface{}              `json:"deactivated"`
	BalanceUpdates            []MetadataBalanceUpdate    `json:"balance_updates"`
	LiquidityBakingEscapeEma  int64                      `json:"liquidity_baking_escape_ema"`
	ImplicitOperationsResults []ImplicitOperationsResult `json:"implicit_operations_results"`
}
