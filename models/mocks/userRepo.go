package mocks

import (
	"context"
	"explorer/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepo is a mock type for auth.UserRepo
type MockUserRepo struct {
	mock.Mock
}

// GetByEmail is mock of UserRepository GetByEmail
func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ret := m.Called(ctx, email)

	var r0 *models.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create is mock of UserRepository Create
func (m *MockUserRepo) CreateUser(ctx context.Context, u *models.User) error {
	ret := m.Called(ctx, u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID is mock of UserRepository GetByID
func (m *MockUserRepo) GetByID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	ret := m.Called(ctx, uid)

	var r0 *models.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}
