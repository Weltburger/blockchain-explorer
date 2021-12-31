package auth

import (
	"context"
	"explorer/models"
	"net/http"

	"github.com/google/uuid"
)

const CtxUserKey = "user"

// UserUsecase defines methods the handler layer expects
// any service it interacts with to implement
type UserUsecase interface {
	SignUp(ctx context.Context, u *models.User) error
	SignIn(ctx context.Context, u *models.User) error
	Get(ctx context.Context, uid uuid.UUID) (*models.User, error)
}

// TokenService defines methods the handler layer expects to interact
// with in regards to producing JWTs as string
type TokenUsecase interface {
	NewPairTokens(ctx context.Context, user string) (*models.TokenDetails, error)
	SavePairTokens(ctx context.Context, td *models.TokenDetails) error
	ValidateAccessToken(ctx context.Context, token string) (*models.ValidationDetails, error)
	ValidateRefreshToken(ctx context.Context, token string) (*models.ValidationDetails, error)
	DeleteTokens(ctx context.Context, r *http.Request) error
	DeleteRefreshToken(ctx context.Context, tokenStr string) error
}
