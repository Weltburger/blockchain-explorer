package models

type TentacledValue struct {
	Int  string     `json:"int,omitempty"`
	Prim PurplePrim `json:"prim,omitempty"`
	Args []Arg10    `json:"args,omitempty"`
}
