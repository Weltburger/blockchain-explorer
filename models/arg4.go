package models

type Arg4 struct {
	String string     `json:"string,omitempty"`
	Prim   PurplePrim `json:"prim,omitempty"`
	Args   []Arg5     `json:"args,omitempty"`
}
