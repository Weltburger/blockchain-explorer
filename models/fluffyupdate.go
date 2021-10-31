package models

type FluffyUpdate struct {
	KeyHash string                `json:"key_hash"`
	Key     Key                   `json:"key"`
	Value   *BigMapDiffValueUnion `json:"value"`
}
