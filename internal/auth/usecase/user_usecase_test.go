package usecase

import (
	"context"
	"explorer/internal/apperrors"
	"explorer/models"
	"explorer/models/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userPW := "strongPassword"
		hashPW, err := hashPassword(userPW)
		assert.NoError(t, err)

		mockUser := &models.User{
			Email:    "1@mail.local",
			Password: userPW,
		}

		mockUserRP := new(mocks.MockUserRepo)
		us := NewUserCase(mockUserRP)

		mockUserRP.On("CreateUser", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*models.User) // arg 0 is context, arg 1 is *User
				userArg.Password = hashPW
			}).Return(nil)

		ctx := context.TODO()
		err = us.SignUp(ctx, mockUser)

		assert.NoError(t, err)

		mockUserRP.AssertExpectations(t)
	})

	t.Run("User already exist error", func(t *testing.T) {
		mockUser := &models.User{
			Email:    "1@mail.local",
			Password: "secretpassword",
		}

		mockUserRP := new(mocks.MockUserRepo)
		us := NewUserCase(mockUserRP)

		mockErr := apperrors.NewConflict("email", mockUser.Email)

		mockUserRP.
			On("CreateUser", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Return(mockErr)

		ctx := context.TODO()
		err := us.SignUp(ctx, mockUser)

		// assert error is error we response with in mock
		assert.EqualError(t, err, mockErr.Error())

		mockUserRP.AssertExpectations(t)
	})
}

func TestSignIn(t *testing.T) {
	// setup valid email/pw combo with hashed password to test method
	// response when provided password is invalid
	email := "1@mail.local"
	validPW := "strongPassword!"
	hashedValidPW, _ := hashPassword(validPW)
	invalidPW := "weakpass"

	mockUserRP := new(mocks.MockUserRepo)
	us := NewUserCase(mockUserRP)

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &models.User{
			Email:    email,
			Password: validPW,
		}

		mockUserResp := &models.User{
			ID:       uid,
			Email:    email,
			Password: hashedValidPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		mockUserRP.
			On("GetByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		err := us.SignIn(ctx, mockUser)

		assert.NoError(t, err)
		mockUserRP.AssertCalled(t, "GetByEmail", mockArgs...)
	})

	t.Run("Invalid email and password combination", func(t *testing.T) {
		mockUser := &models.User{
			Email:    email,
			Password: invalidPW,
		}

		mockUserResp := &models.User{
			Email:    email,
			Password: hashedValidPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		mockUserRP.On("GetByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		err := us.SignIn(ctx, mockUser)

		waitError := apperrors.NewInternal()

		assert.Error(t, err)
		assert.EqualError(t, err, waitError.Error())
		mockUserRP.AssertCalled(t, "GetByEmail", mockArgs...)
	})
}
