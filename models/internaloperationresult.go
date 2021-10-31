package models

type InternalOperationResult struct {
	Kind        ImplicitOperationsResultKind       `json:"kind"`
	Source      Source                             `json:"source"`
	Nonce       int64                              `json:"nonce"`
	Amount      string                             `json:"amount"`
	Destination string                             `json:"destination"`
	Parameters  *InternalOperationResultParameters `json:"parameters,omitempty"`
	Result      Result                             `json:"result"`
}
