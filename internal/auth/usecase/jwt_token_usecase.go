package usecase

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
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
func (s *tokenUsecase) NewPairTokens(ctx context.Context, uid string) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}

	currentUTime := time.Now().Unix()
	td.AtExpires = currentUTime + s.AccessExpirationSecs
	td.AccessUuid = uuid.NewString()

	claimsAt := accessTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			Id:        td.AccessUuid,
			IssuedAt:  currentUTime,
			ExpiresAt: td.AtExpires,
			Issuer:    "Block-explorer",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsAt)
	signedAccessToken, err := accessToken.SignedString(s.PrivKey)
	if err != nil {
		log.Printf("Failed to sign id token: %v\n", err)
		return nil, err
	}

	td.AccessToken = signedAccessToken

	td.RtExpires = currentUTime + s.RefreshExpirationSecs
	td.RefreshUuid = fmt.Sprintf("%s:%s", td.AccessUuid, uid)

	claimsRt := refreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			Id:        td.RefreshUuid,
			IssuedAt:  currentUTime,
			ExpiresAt: td.RtExpires,
			Issuer:    "Block-explorer",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRt)
	signedRefreshToken, err := refreshToken.SignedString([]byte(s.RefreshSecret))
	if err != nil {
		log.Printf("Failed to sign refresh token: %v\n", err)
		return nil, err
	}

	td.RefreshToken = signedRefreshToken

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
func (s *tokenUsecase) ValidateAccessToken(ctx context.Context, tokenStr string) (*models.ValidationDetails, error) {

	claims, err := validAcToken(tokenStr, s.PubKey) // uses public RSA key
	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse accessToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization(err.Error())
	}

	userId, err := s.TokenRepository.FetchAuth(ctx, claims.Id)
	if err != nil {
		log.Printf("Error fetch tokenId %s from Redis: %v", claims.Id, err)
		return nil, apperrors.NewInternal()
	}

	log.Println(userId)
	log.Println(claims.UID)

	if userId != claims.UID {
		return nil, apperrors.NewAuthorization("Error compare userId in DB and token.")
	}

	return &models.ValidationDetails{
		TokenUuid: claims.Id,
		UserId:    claims.UID,
	}, nil
}

func (s *tokenUsecase) ValidateRefreshToken(ctx context.Context, tokenStr string) (*models.ValidationDetails, error) {

	refreshClaims, err := validRfToken(tokenStr, s.RefreshSecret)
	if err != nil {
		log.Printf("Unable to validate or parse refreshToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization(err.Error())
	}

	return &models.ValidationDetails{
		TokenUuid: refreshClaims.Id,
		UserId:    refreshClaims.UID,
	}, nil
}

// DeleteTokens handler delete tokens from DB
func (s *tokenUsecase) DeleteTokens(ctx context.Context, r *http.Request) error {
	tokenStr := extractToken(r)
	if tokenStr == "" {
		return apperrors.NewBadRequest("Error extract token from request header.")
	}
	accessClaims, err := validAcToken(tokenStr, s.PubKey)
	if err != nil {
		return apperrors.NewAuthorization(err.Error())
	}

	delAcError := s.TokenRepository.DeleteAccessToken(ctx, accessClaims.UID)
	if delAcError != nil {
		return delAcError
	}

	refKey := fmt.Sprintf("%s:%s", accessClaims.Id, accessClaims.UID)
	delRfError := s.TokenRepository.DeleteRefreshToken(ctx, refKey)
	if delRfError != nil {
		return delRfError
	}

	return nil
}

// DeleteTokens handler delete tokens from DB
func (s *tokenUsecase) DeleteRefreshToken(ctx context.Context, tokenStr string) error {

	err := s.TokenRepository.DeleteRefreshToken(ctx, tokenStr)
	if err != nil {
		return err
	}

	return nil
}
