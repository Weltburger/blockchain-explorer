package models

type Arg1 struct {
	Prim PurplePrim `json:"prim,omitempty"`
	Args []Arg2     `json:"args,omitempty"`
	Int  string     `json:"int,omitempty"`
}
