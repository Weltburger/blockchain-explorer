package models

type ContentParameters struct {
	Entrypoint string          `json:"entrypoint"`
	Value      *HilariousValue `json:"value"`
}
