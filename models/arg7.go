package models

type Arg7 struct {
	String string            `json:"string,omitempty"`
	Int    string            `json:"int,omitempty"`
	Prim   PurplePrim        `json:"prim,omitempty"`
	Args   []StorageArgClass `json:"args,omitempty"`
}
