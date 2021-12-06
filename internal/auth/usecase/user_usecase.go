package usecase

import (
	"context"
	"time"

	"explorer/internal/apperrors"
	"explorer/internal/auth"
	"explorer/models"

	"github.com/google/uuid"
)

type UserCase struct {
	userRepo auth.UserRepo
}

func NewUserUseCase(userRepo auth.UserRepo) auth.UserUsecase {
	return &UserCase{
		userRepo: userRepo,
	}
}

func (u *UserCase) SignUp(ctx context.Context, usr *models.User) error {
	// hash user password
	pwd, err := hashPassword(usr.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     usr.Password,
		Password:  pwd,
		CreatedAt: time.Now(),
	}

	// store new user to DB
	return u.userRepo.CreateUser(ctx, user)
}

func (u *UserCase) SignIn(ctx context.Context, usr *models.User) error {
	user, err := u.userRepo.GetByEmail(ctx, usr.Email)
	if err != nil {
		return apperrors.NewNotFound("username", usr.Email)
	}

	// verify password
	match := doPasswordsMatch(user.Password, usr.Password)
	if !match {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	return nil

}
