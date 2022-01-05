package usecase

import (
	"context"
	"log"

	"explorer/internal/apperrors"
	"explorer/internal/auth"
	"explorer/models"
)

type UserCase struct {
	userRepo auth.UserRepo
}

func NewUserCase(userRepo auth.UserRepo) auth.UserUsecase {
	return &UserCase{
		userRepo: userRepo,
	}
}

// Signup handler use for register new user
func (u *UserCase) SignUp(ctx context.Context, usr *models.User) error {
	// hash user password
	pwd, err := hashPassword(usr.Password)
	if err != nil {
		log.Printf("Error to hash password for user: %v\n", usr.Email)
		return apperrors.NewInternal()
	}

	// change passsword to hash
	usr.Password = pwd

	// store new user to DB
	return u.userRepo.CreateUser(ctx, usr)
}

// Signin handler use for fetch user from DB and compare passwords
func (u *UserCase) SignIn(ctx context.Context, usr *models.User) error {
	userFetched, err := u.userRepo.GetByEmail(ctx, usr.Email)
	if err != nil {
		return err
	}

	// verify password
	match, err := doPasswordsMatch(userFetched.Password, usr.Password)
	if err != nil {
		log.Printf("Error compare hash: %v\n", err)
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	// add autogenerated fields
	*usr = *userFetched

	return nil

}
