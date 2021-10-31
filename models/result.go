package models

type Result struct {
	Status              Status                                  `json:"status"`
	Storage             *ResultStorage                          `json:"storage"`
	BigMapDiff          []ResultBigMapDiff                      `json:"big_map_diff,omitempty"`
	BalanceUpdates      []ImplicitOperationsResultBalanceUpdate `json:"balance_updates,omitempty"`
	ConsumedGas         string                                  `json:"consumed_gas"`
	ConsumedMilligas    string                                  `json:"consumed_milligas"`
	StorageSize         string                                  `json:"storage_size,omitempty"`
	PaidStorageSizeDiff string                                  `json:"paid_storage_size_diff,omitempty"`
	LazyStorageDiff     []ResultLazyStorageDiff                 `json:"lazy_storage_diff,omitempty"`
}
