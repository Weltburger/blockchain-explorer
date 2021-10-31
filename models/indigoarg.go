package models

type IndigoArg struct {
	Prim string        `json:"prim"`
	Args []IndecentArg `json:"args"`
}
