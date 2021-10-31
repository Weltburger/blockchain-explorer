package models

type OperationResult struct {
	Status              Status                                  `json:"status"`
	Storage             *OperationResultStorage                 `json:"storage"`
	BigMapDiff          []OperationResultBigMapDiff             `json:"big_map_diff,omitempty"`
	ConsumedGas         string                                  `json:"consumed_gas"`
	ConsumedMilligas    string                                  `json:"consumed_milligas"`
	StorageSize         string                                  `json:"storage_size,omitempty"`
	LazyStorageDiff     []OperationResultLazyStorageDiff        `json:"lazy_storage_diff,omitempty"`
	BalanceUpdates      []ImplicitOperationsResultBalanceUpdate `json:"balance_updates,omitempty"`
	PaidStorageSizeDiff string                                  `json:"paid_storage_size_diff,omitempty"`
}
