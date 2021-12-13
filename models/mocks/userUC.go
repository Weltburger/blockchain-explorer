package mocks

import (
	"context"
	"explorer/models"

	"github.com/stretchr/testify/mock"
)

// MockUserCase is a mock type for UserUseCase
type MockUserCase struct {
	mock.Mock
}

// SignUp is mock for User Usecase SignUp
func (m *MockUserCase) SignUp(ctx context.Context, user *models.User) error {
	// args that will be passed to "Return" in the tests, when function is called.
	ret := m.Called(ctx, user)

	//value passed to "Return"
	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Error(0)
	}

	return r0
}

// SignIn is mock for User Usecase SignIn
func (m *MockUserCase) SignIn(ctx context.Context, user *models.User) error {
	// args that will be passed to "Return" in the tests, when function is called.
	ret := m.Called(ctx, user)

	//value passed to "Return"
	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Error(0)
	}

	return r0
}
