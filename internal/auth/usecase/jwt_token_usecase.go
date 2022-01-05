package usecase

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"explorer/internal/apperrors"
	"explorer/internal/auth"
	"explorer/models"
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
func (s *tokenUsecase) GenerateTokens(ctx context.Context, uid string) (*models.TokenPair, error) {
	tp := &models.TokenPair{}

	// make access token
	accessTokenData, err := makeAccessT(uid, s.PrivKey, s.AccessExpirationSecs)
	if err != nil {
		return nil, err
	}

	// make refresh token
	refreshTokenData, err := makeRefreshT(strings.Join([]string{accessTokenData.TokenId, uid}, ":"), []byte(s.RefreshSecret), s.RefreshExpirationSecs)
	if err != nil {
		return nil, err
	}

	at := time.Unix(accessTokenData.ExpiredAt, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(refreshTokenData.ExpiredAt, 0)
	now := time.Now()

	// set access token to valid list
	if err := s.TokenRepository.SetToken(ctx, accessTokenData.TokenId, accessTokenData.Token, at.Sub(now)); err != nil {
		return tp, err
	}

	// set access token to valid list
	if err := s.TokenRepository.SetToken(ctx, refreshTokenData.TokenId, refreshTokenData.Token, rt.Sub(now)); err != nil {
		return tp, err
	}

	tp.AccessToken = accessTokenData.Token
	tp.RefreshToken = refreshTokenData.Token

	return tp, nil
}

// DeleteTokens handler delete tokens from DB
func (s *tokenUsecase) DeleteTokens(ctx context.Context, r *http.Request) error {
	tokenStr := extractToken(r)
	if tokenStr == "" {
		return apperrors.NewBadRequest("Error extract token from request header.")
	}
	tokenClaims, err := validAccessT(tokenStr, s.PubKey)
	if err != nil {
		return apperrors.NewAuthorization(err.Error())
	}

	delAcError := s.TokenRepository.DeleteToken(ctx, tokenClaims.TokenId)
	if delAcError != nil {
		return delAcError
	}

	refKey := fmt.Sprintf("%s:%s", tokenClaims.TokenId, tokenClaims.UserId)
	delRfError := s.TokenRepository.DeleteToken(ctx, refKey)
	if delRfError != nil {
		return delRfError
	}

	return nil
}

// ValidateToken validates the id token jwt string
// It returns the user extract from the TokenCustomClaims
func (s *tokenUsecase) ValidateToken(ctx context.Context, tokenStr string, isRefresh bool) (*models.ValidationDetails, error) {

	customClaims := &TokenCustomClaims{}
	if isRefresh {
		claims, err := validRefreshT(tokenStr, []byte(s.RefreshSecret))
		if err != nil {
			log.Printf("Unable to validate or parse accessToken - Error: %v\n", err)
			return nil, apperrors.NewAuthorization(err.Error())
		}
		*customClaims = *claims
	} else {
		claims, err := validAccessT(tokenStr, s.PubKey) // uses public RSA key
		if err != nil {
			log.Printf("Unable to validate or parse accessToken - Error: %v\n", err)
			return nil, apperrors.NewAuthorization(err.Error())
		}
		*customClaims = *claims
	}

	tokenDB, err := s.TokenRepository.FetchToken(ctx, customClaims.TokenId)
	if err != nil {
		log.Printf("Error fetch token %s from Redis: %v", customClaims.TokenId, err)
		return nil, apperrors.NewInternal()
	}

	if tokenStr != tokenDB {
		return nil, apperrors.NewAuthorization("Error compare user token with DB token.")
	}

	return &models.ValidationDetails{
		TokenId: customClaims.TokenId,
		UserId:  customClaims.UserId,
	}, nil
}

// DeleteRefreshTokens handler delete refresh token from DB
func (s *tokenUsecase) DeleteRefreshToken(ctx context.Context, tokenStr string) error {

	err := s.TokenRepository.DeleteToken(ctx, tokenStr)
	if err != nil {
		return err
	}

	return nil
}
