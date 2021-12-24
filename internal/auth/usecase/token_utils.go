package usecase

import (
	"crypto/rsa"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// UserToToken struct holds data which will be insert to idToken
type userToToken struct {
	UID   string
	Email string
}

// idTokenCustomClaims holds structure of jwt claims of idToken
type accessTokenCustomClaims struct {
	User *userToToken `json:"user"`
	jwt.StandardClaims
}

type refreshTokenCustomClaims struct {
	UID string `json:"uid"`
	jwt.StandardClaims
}

// generateIDToken generates an IDToken which is a jwt with myCustomClaims
// Could call this GenerateIDTokenString, but the signature makes this fairly clear
// func generateAccessToken(u *models.User, key *rsa.PrivateKey, exp int64) (string, error) {
// 	unixTime := time.Now().Unix()
// 	tokenExp := unixTime + exp
// 	accessTokenUuid := uuid.NewString()

// 	claims := accessTokenCustomClaims{
// 		User: &userToToken{UID: u.ID.String(),
// 			Email: u.Email},
// 		StandardClaims: jwt.StandardClaims{
// 			Id:        accessTokenUuid,
// 			IssuedAt:  unixTime,
// 			ExpiresAt: tokenExp,
// 		},
// 	}

// 	idToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 	signedToken, err := idToken.SignedString(key)
// 	if err != nil {
// 		log.Println("Failed to sign id token")
// 		return "", err
// 	}

// 	return signedToken, nil
// }

// refreshTokenData holds the actual signed jwt string along with the ID
// We return the id so it can be used without re-parsing the JWT from signed string
// type refreshTokenData struct {
// 	SS        string
// 	ID        uuid.UUID
// 	ExpiresIn time.Duration
// }

// refreshTokenCustomClaims holds the payload of a refresh token
// This can be used to extract user id for subsequent
// application operations (IE, fetch user in Redis)

// generateRefreshToken creates a refresh token
// The refresh token stores only the user's ID, a string
// func generateRefreshToken(uid uuid.UUID, key string, exp int64) (*refreshTokenData, error) {
// 	currentTime := time.Now()
// 	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)
// 	tokenID, err := uuid.NewRandom() // v4 uuid in the google uuid lib
// 	if err != nil {
// 		log.Println("Failed to generate refresh token ID")
// 		return nil, err
// 	}

// 	claims := refreshTokenCustomClaims{
// 		UID: uid,
// 		StandardClaims: jwt.StandardClaims{
// 			IssuedAt:  currentTime.Unix(),
// 			ExpiresAt: tokenExp.Unix(),
// 			Id:        tokenID.String(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	ss, err := token.SignedString([]byte(key))

// 	if err != nil {
// 		log.Println("Failed to sign refresh token string")
// 		return nil, err
// 	}

// 	return &refreshTokenData{
// 		SS:        ss,
// 		ID:        tokenID,
// 		ExpiresIn: tokenExp.Sub(currentTime),
// 	}, nil
// }

// validateIDToken returns the token's claims if the token is valid
func validateIDToken(tokenString string, key *rsa.PublicKey) (*accessTokenCustomClaims, error) {
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
