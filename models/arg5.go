package models

type Arg5 struct {
	String string     `json:"string,omitempty"`
	Prim   PurplePrim `json:"prim,omitempty"`
	Args   []Arg6     `json:"args,omitempty"`
}
