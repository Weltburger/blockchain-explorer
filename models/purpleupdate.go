package models

type PurpleUpdate struct {
	KeyHash string         `json:"key_hash"`
	Key     KeyElement     `json:"key"`
	Value   TentacledValue `json:"value"`
}
