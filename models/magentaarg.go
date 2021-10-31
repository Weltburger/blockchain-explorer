package models

type MagentaArg struct {
	Prim  PurplePrim `json:"prim,omitempty"`
	Args  []Arg13    `json:"args,omitempty"`
	Int   string     `json:"int,omitempty"`
	Bytes string     `json:"bytes,omitempty"`
}
