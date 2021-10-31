package models

type MetadataBalanceUpdate struct {
	Kind     BalanceUpdateKind `json:"kind"`
	Contract string            `json:"contract,omitempty"`
	Change   string            `json:"change"`
	Origin   Origin            `json:"origin"`
	Category Category          `json:"category,omitempty"`
	Delegate string            `json:"delegate,omitempty"`
	Cycle    int64             `json:"cycle,omitempty"`
}
