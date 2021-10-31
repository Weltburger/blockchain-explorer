package models

type Key struct {
	Int   string       `json:"int,omitempty"`
	Prim  FluffyPrim   `json:"prim,omitempty"`
	Args  []KeyElement `json:"args,omitempty"`
	Bytes Bytes        `json:"bytes,omitempty"`
}
