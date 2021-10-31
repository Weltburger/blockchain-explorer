package models

type FluffyValue struct {
	Prim PurplePrim  `json:"prim"`
	Args []StickyArg `json:"args"`
}
