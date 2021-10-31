package models

type PurpleDiff struct {
	Action  Action         `json:"action"`
	Updates []PurpleUpdate `json:"updates"`
}
