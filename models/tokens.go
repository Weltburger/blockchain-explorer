package models

import "github.com/google/uuid"

// RefreshToken stores token properties that
// are accessed in multiple application layers
type RefreshToken struct {
	ID  uuid.UUID `json:"-"`
	UID uuid.UUID `json:"-"`
	SS  string    `json:"refreshToken"`
}

// IDToken stores token properties that
// are accessed in multiple application layers
type IDToken struct {
	SS string `json:"idToken"`
}

// TokenDetails used for pair of id and refresh tokens with addition parameters
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
