package models

type InternalOperationResultParameters struct {
	Entrypoint Entrypoint     `json:"entrypoint"`
	Value      *IndecentValue `json:"value"`
}
