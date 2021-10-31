package models

type KeyElement struct {
	Prim  PurplePrim                        `json:"prim,omitempty"`
	Args  []ImplicitOperationsResultStorage `json:"args,omitempty"`
	Int   string                            `json:"int,omitempty"`
	Bytes string                            `json:"bytes,omitempty"`
}
