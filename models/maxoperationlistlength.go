package models

type MaxOperationListLength struct {
	MaxSize int64 `json:"max_size"`
	MaxOp   int64 `json:"max_op,omitempty"`
}
