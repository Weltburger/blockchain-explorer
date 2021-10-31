package models

type FriskyArg struct {
	Prim PurplePrim       `json:"prim,omitempty"`
	Args []MischievousArg `json:"args,omitempty"`
	Int  string           `json:"int,omitempty"`
}
