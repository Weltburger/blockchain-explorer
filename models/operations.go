package models

type Operations struct {
	Kind  OperationsKind `json:"kind"`
	Level int64          `json:"level"`
}
