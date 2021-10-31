package models

type BraggadociousArg struct {
	Bytes string         `json:"bytes,omitempty"`
	Prim  string         `json:"prim,omitempty"`
	Args  []TentacledArg `json:"args,omitempty"`
	Int   string         `json:"int,omitempty"`
}
