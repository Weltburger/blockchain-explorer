package models

type StorageArg struct {
	Prim FluffyPrim        `json:"prim,omitempty"`
	Args []StorageArgClass `json:"args,omitempty"`
	Int  string            `json:"int,omitempty"`
}
