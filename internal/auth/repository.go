package auth

import (
	"context"
	"explorer/models"

	"github.com/google/uuid"
)

// UserRepo defines methods the service layer expects
// any repository it interacts with to implement
type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, uid uuid.UUID) (*models.User, error)
}
