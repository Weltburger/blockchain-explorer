package models

type IndigoValue struct {
	Prim PurplePrim `json:"prim,omitempty"`
	Args []Arg7     `json:"args,omitempty"`
	Int  string     `json:"int,omitempty"`
}
