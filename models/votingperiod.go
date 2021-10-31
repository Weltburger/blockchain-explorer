package models

type VotingPeriod struct {
	Index         int64  `json:"index"`
	Kind          string `json:"kind"`
	StartPosition int64  `json:"start_position"`
}
