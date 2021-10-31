package models

type Arg2 struct {
	Bytes string `json:"bytes,omitempty"`
	Prim  string `json:"prim,omitempty"`
	Args  []Arg3 `json:"args,omitempty"`
	Int   string `json:"int,omitempty"`
}
