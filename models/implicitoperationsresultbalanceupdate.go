package models

type ImplicitOperationsResultBalanceUpdate struct {
	Kind     BalanceUpdateKind `json:"kind"`
	Contract string            `json:"contract"`
	Change   string            `json:"change"`
	Origin   Origin            `json:"origin"`
}
