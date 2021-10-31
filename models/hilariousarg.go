package models

type HilariousArg struct {
	Prim PurplePrim     `json:"prim,omitempty"`
	Args []AmbitiousArg `json:"args,omitempty"`
	Int  string         `json:"int,omitempty"`
}
