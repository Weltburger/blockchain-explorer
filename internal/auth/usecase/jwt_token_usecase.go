package usecase

import (
	"context"
	"crypto/rsa"
	"log"
	"time"

	"explorer/internal/apperrors"
	"explorer/internal/auth"
	"explorer/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// tokenUsecase used for injecting an implementation of TokenRepository
// for use in service methods along with keys and secrets for
// signing JWTs
type tokenUsecase struct {
	TokenRepository       auth.TokenRepo
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	AccessExpirationSecs  int64
	RefreshExpirationSecs int64
}

// TSConfig will hold repositories that will eventually be injected into this
// this service layer
type TSConfig struct {
	TokenRepository       auth.TokenRepo
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	AccessExpirationSecs  int64
	RefreshExpirationSecs int64
}

// NewTokenUsecase is a factory function for
// initializing a TokenUsecase with its repository layer dependencies
func NewTokenUsecase(c *TSConfig) auth.TokenUsecase {
	return &tokenUsecase{
		TokenRepository:       c.TokenRepository,
		PrivKey:               c.PrivKey,
		PubKey:                c.PubKey,
		RefreshSecret:         c.RefreshSecret,
		AccessExpirationSecs:  c.AccessExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

// NewPairTokens creates fresh id and refresh tokens for the current user
func (s *tokenUsecase) NewPairTokens(ctx context.Context, u *models.User) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}

	currentUTime := time.Now().Unix()
	td.AtExpires = currentUTime + s.AccessExpirationSecs
	td.AccessUuid = uuid.NewString()

	claimsAt := accessTokenCustomClaims{
		User: &userToToken{UID: u.ID.String(),
			Email: u.Email},
		StandardClaims: jwt.StandardClaims{
			Id:        td.AccessUuid,
			IssuedAt:  currentUTime,
			ExpiresAt: td.AtExpires,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsAt)
	signedAccessToken, err := accessToken.SignedString(s.PrivKey)
	if err != nil {
		log.Printf("Failed to sign id token: %v\n", err)
		return nil, apperrors.NewInternal()
	}

	td.AccessToken = signedAccessToken

	td.RtExpires = currentUTime + s.RefreshExpirationSecs
	td.RefreshUuid = uuid.NewString()

	claimsRt := refreshTokenCustomClaims{
		UID: u.ID.String(),
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentUTime,
			ExpiresAt: td.RtExpires,
			Id:        td.RefreshUuid,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRt)
	signedRefreshToken, err := refreshToken.SignedString([]byte(s.RefreshSecret))
	if err != nil {
		log.Printf("Failed to sign refresh token: %v\n", err)
		return nil, err
	}

	td.RefreshToken = signedRefreshToken

	// refreshToken, err := generateRefreshToken(u.ID, s.RefreshSecret, s.RefreshExpirationSecs)
	// if err != nil {
	// 	log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.ID, err.Error())
	// 	return nil, apperrors.NewInternal()
	// }

	// delete user's current refresh token (used when refreshing idToken)
	// if prevTokenID != "" {
	// 	if err := s.TokenRepository.DeleteRefreshToken(ctx, u.ID.String(), prevTokenID); err != nil {
	// 		log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.ID.String(), prevTokenID)
	// 	}
	// }

	return td, nil
}

// SavePairTokens save tokens to Redis DB
func (s *tokenUsecase) SavePairTokens(ctx context.Context, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	// set freshly minted access token to valid list
	if err := s.TokenRepository.SetAccessToken(ctx, td.AccessUuid, td.AccessToken, at.Sub(now)); err != nil {
		return err
	}

	// set freshly minted refresh token to valid list
	if err := s.TokenRepository.SetRefreshToken(ctx, td.RefreshUuid, td.RefreshToken, rt.Sub(now)); err != nil {
		return err
	}

	return nil

}

// ValidateIDToken validates the id token jwt string
// It returns the user extract from the IDTokenCustomClaims
func (s *tokenUsecase) ValidateIDToken(ctx context.Context, tokenString string) (*models.User, error) {
	claims, err := validateIDToken(tokenString, s.PubKey) // uses public RSA key

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization("Unable to verify user from idToken")
	}

	return &models.User{
		ID:    uuid.MustParse(claims.User.UID),
		Email: claims.User.Email,
	}, nil
}

func (s *tokenUsecase) RefreshToken(ctx context.Context) error {
	return nil
}
