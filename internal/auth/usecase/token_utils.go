package usecase

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// idTokenCustomClaims holds structure of jwt claims of idToken
type accessTokenCustomClaims struct {
	UID string `json:"uid"`
	jwt.StandardClaims
}

type refreshTokenCustomClaims struct {
	UID string `json:"uid"`
	jwt.StandardClaims
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

// validAcToken returns the token's claims if the token is valid
func validAcToken(tokenString string, key *rsa.PublicKey) (*accessTokenCustomClaims, error) {
	claims := &accessTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*accessTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}

// validRfToken returns the token's claims if the token is valid
func validRfToken(tokenString string, key string) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*refreshTokenCustomClaims)
	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}
