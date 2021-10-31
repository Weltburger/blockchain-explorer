package models

type ResultBigMapDiff struct {
	Action  Action         `json:"action"`
	BigMap  string         `json:"big_map"`
	KeyHash string         `json:"key_hash"`
	Key     KeyElement     `json:"key"`
	Value   TentacledValue `json:"value"`
}
