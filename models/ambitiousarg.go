package models

type AmbitiousArg struct {
	Prim string       `json:"prim,omitempty"`
	Args []CunningArg `json:"args,omitempty"`
	Int  string       `json:"int,omitempty"`
}
