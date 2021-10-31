package models

type StorageArgClass struct {
	Bytes  string      `json:"bytes,omitempty"`
	Prim   PurplePrim  `json:"prim,omitempty"`
	Args   []FluffyArg `json:"args,omitempty"`
	Int    string      `json:"int,omitempty"`
	String string      `json:"string,omitempty"`
}
