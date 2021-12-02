package usecase

import (
	"context"
	"fmt"
	"time"

	"explorer/internal/auth"
	"explorer/internal/auth/apperrors"
	"explorer/models"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

type tokenUser struct {
	ID    uuid.UUID `json:"uid"`
	Email string    `json:"email"`
}
type AuthClaims struct {
	jwt.StandardClaims
	User tokenUser `json:"user"`
}

type AuthUseCase struct {
	userRepo       auth.UserRepo
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(
	userRepo auth.UserRepo,
	signingKey []byte,
	tokenTTLSeconds time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, username, password string) error {
	pwd, err := hashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     username,
		Password:  pwd,
		CreatedAt: time.Now(),
	}

	return a.userRepo.CreateUser(ctx, user)
}

func (a *AuthUseCase) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", apperrors.NewNotFound("username", user.Email)
	}

	// verify password
	match := doPasswordsMatch(user.Password, password)
	if !match {
		return "", apperrors.NewAuthorization("Invalid email and password combination")
	}

	expiresTime := time.Now().Add(a.expireDuration)
	claims := AuthClaims{
		User: tokenUser{
			ID:    user.ID,
			Email: user.Email,
		},
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  jwt.At(time.Now()),
			ExpiresAt: jwt.At(expiresTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.signingKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return &models.User{
			ID:    claims.User.ID,
			Email: claims.User.Email,
		}, nil
	}

	return nil, apperrors.NewAuthorization("Invalid token")
}
