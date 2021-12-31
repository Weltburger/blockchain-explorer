package models

// ValidationDetails user and token id pair
type ValidationDetails struct {
	TokenUuid string
	UserId    string
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

// TokenPair send send consumer after signin and refresh
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
