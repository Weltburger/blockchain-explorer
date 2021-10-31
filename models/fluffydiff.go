package models

type FluffyDiff struct {
	Action  Action         `json:"action"`
	Updates []FluffyUpdate `json:"updates"`
}
