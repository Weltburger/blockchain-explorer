package models

type MischievousArg struct {
	Prim PurplePrim         `json:"prim,omitempty"`
	Args []BraggadociousArg `json:"args,omitempty"`
	Int  string             `json:"int,omitempty"`
}
