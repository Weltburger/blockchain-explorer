package models

type EndorsementClass struct {
	Branch     Predecessor `json:"branch"`
	Operations Operations  `json:"operations"`
	Signature  string      `json:"signature"`
}
