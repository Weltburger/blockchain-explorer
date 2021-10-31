package models

type Content struct {
	Kind         ImplicitOperationsResultKind `json:"kind"`
	Endorsement  *EndorsementClass            `json:"endorsement,omitempty"`
	Slot         int64                        `json:"slot,omitempty"`
	Metadata     ContentMetadata              `json:"metadata"`
	Source       string                       `json:"source,omitempty"`
	Fee          string                       `json:"fee,omitempty"`
	Counter      string                       `json:"counter,omitempty"`
	GasLimit     string                       `json:"gas_limit,omitempty"`
	StorageLimit string                       `json:"storage_limit,omitempty"`
	Amount       string                       `json:"amount,omitempty"`
	Destination  string                       `json:"destination,omitempty"`
	Parameters   *ContentParameters           `json:"parameters,omitempty"`
}
