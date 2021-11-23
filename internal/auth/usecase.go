package auth

import (
	"context"
	"explorer/models"
)

const CtxUserKey = "user"

// UserUsecase defines methods the handler layer expects
// any service it interacts with to implement
type UserUsecase interface {
	SignUp(ctx context.Context, email, password string) error
	SignIn(ctx context.Context, email, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (models.User, error)
}
