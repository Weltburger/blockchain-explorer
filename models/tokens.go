package models

// ValidationDetails user and token id pair
type ValidationDetails struct {
	UserId  string
	TokenId string
}

// TokenPair send send consumer after signin and refresh
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
