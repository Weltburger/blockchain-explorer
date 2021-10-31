package models

type CunningArg struct {
	Bytes string     `json:"bytes,omitempty"`
	Prim  FluffyPrim `json:"prim,omitempty"`
}
