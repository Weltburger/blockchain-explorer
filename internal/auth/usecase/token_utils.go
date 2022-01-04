package usecase

import (
	"crypto/rsa"
	"explorer/internal/apperrors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// TokenCustomClaims holds structure of jwt claims
type TokenCustomClaims struct {
	UserId  string `json:"user_id"`
	TokenId string `json:"token_id"`
	jwt.StandardClaims
}

// TokenData struct holds all generated token metadata
type tokenData struct {
	TokenId   string
	ExpiredAt int64
	Token     string
}

//get the token from the request header
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	tokenArr := strings.Split(bearerToken, " ")
	if len(tokenArr) == 2 && strings.Contains(tokenArr[0], "Bearer") {
		return tokenArr[1]
	}
	return ""
}

// makeAccessT function create access token
func makeAccessT(id string, key *rsa.PrivateKey, tokenDuration int64) (*tokenData, error) {
	td := &tokenData{}
	currentUTime := time.Now().Unix()
	td.ExpiredAt = currentUTime + tokenDuration
	td.TokenId = uuid.NewString()

	customClaim := TokenCustomClaims{
		UserId:  id,
		TokenId: td.TokenId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentUTime,
			ExpiresAt: td.ExpiredAt,
			Issuer:    "Block-explorer",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, customClaim)
	signedAccessToken, err := accessToken.SignedString(key)
	if err != nil {
		log.Printf("Failed to sign access token: %v\n", err)
		return nil, apperrors.NewInternal()
	}

	td.Token = signedAccessToken

	return td, nil
}

// makeRefreshT function crate refresh token
func makeRefreshT(tokenId string, key []byte, tokenDuration int64) (*tokenData, error) {
	td := &tokenData{}

	currentUTime := time.Now().Unix()
	td.ExpiredAt = currentUTime + tokenDuration
	td.TokenId = tokenId

	tokenAndUserId := strings.Split(tokenId, ":")
	if len(tokenAndUserId) != 2 {
		log.Printf("Error split token and user id: %s\n", tokenId)
		return nil, apperrors.NewInternal()
	}

	customClaim := TokenCustomClaims{
		UserId:  tokenAndUserId[1],
		TokenId: td.TokenId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentUTime,
			ExpiresAt: td.ExpiredAt,
			Issuer:    "Block-explorer",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaim)
	signedRefreshToken, err := refreshToken.SignedString(key)
	if err != nil {
		log.Printf("Failed to sign token: %v\n", err)
		return nil, apperrors.NewInternal()
	}

	td.Token = signedRefreshToken

	return td, nil
}

// validAccessT returns the token's claims if the token is valid
func validAccessT(tokenString string, key *rsa.PublicKey) (*TokenCustomClaims, error) {

	claims := &TokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Token is invalid")
	}

	claims, ok := token.Claims.(*TokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("Token valid but couldn't parse claims")
	}

	return claims, nil
}

// validAccessT returns the token's claims if the token is valid
func validRefreshT(tokenString string, key []byte) (*TokenCustomClaims, error) {

	claims := &TokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is invalid")
	}

	claims, ok := token.Claims.(*TokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("Token valid but couldn't parse claims")
	}

	return claims, nil
}
