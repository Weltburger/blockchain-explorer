package models

type OperationResultBigMapDiff struct {
	Action  Action                `json:"action"`
	BigMap  string                `json:"big_map"`
	KeyHash string                `json:"key_hash"`
	Key     Key                   `json:"key"`
	Value   *BigMapDiffValueUnion `json:"value"`
}
