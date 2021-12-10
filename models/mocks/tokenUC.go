package mocks

import (
	"context"
	"explorer/models"

	"github.com/stretchr/testify/mock"
)

// MockTokenUsecase is a mock type for TokenUsecase
type MockTokenCase struct {
	mock.Mock
}

// NewPairFromUser mocks concrete NewPairFromUser
func (m *MockTokenCase) NewPairTokens(ctx context.Context, u *models.User, prevTokenID string) (*models.TokenPair, error) {
	ret := m.Called(ctx, u, prevTokenID)

	// first value passed to "Return"
	var r0 *models.TokenPair
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.TokenPair)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateIDToken mocks concrete ValidateIDToken
func (m *MockTokenCase) ValidateIDToken(ctx context.Context, tokenString string) (*models.User, error) {
	ret := m.Called(tokenString)

	// first value passed to "Return"
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
