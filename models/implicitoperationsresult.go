package models

type ImplicitOperationsResult struct {
	Kind             ImplicitOperationsResultKind            `json:"kind"`
	Storage          []ImplicitOperationsResultStorage       `json:"storage"`
	BalanceUpdates   []ImplicitOperationsResultBalanceUpdate `json:"balance_updates"`
	ConsumedGas      string                                  `json:"consumed_gas"`
	ConsumedMilligas string                                  `json:"consumed_milligas"`
	StorageSize      string                                  `json:"storage_size"`
}
