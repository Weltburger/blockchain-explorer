package models

type VotingPeriodInfo struct {
	VotingPeriod VotingPeriod `json:"voting_period"`
	Position     int64        `json:"position"`
	Remaining    int64        `json:"remaining"`
}
