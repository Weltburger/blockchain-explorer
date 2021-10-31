package models

type PurpleArg struct {
	Prim PurplePrim        `json:"prim"`
	Args []StorageArgClass `json:"args"`
}
