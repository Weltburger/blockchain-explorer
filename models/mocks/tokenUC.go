package mocks

import (
	"context"
	"explorer/models"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockTokenUsecase is a mock type for TokenUsecase
type MockTokenCase struct {
	mock.Mock
}

// NewPairFromUser mocks concrete NewPairFromUser
func (m *MockTokenCase) GenerateTokens(ctx context.Context, user string) (*models.TokenPair, error) {
	ret := m.Called(ctx, user)

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

// DeleteTokens mocks DeleteTokens method
func (m *MockTokenCase) DeleteTokens(ctx context.Context, r *http.Request) error {
	ret := m.Called(ctx, r)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateToken mocks concrete ValidateToken
func (m *MockTokenCase) ValidateToken(ctx context.Context, token string, refresh bool) (*models.ValidationDetails, error) {
	ret := m.Called(ctx, token, refresh)

	// first value passed to "Return"
	var r0 *models.ValidationDetails
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.ValidationDetails)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRefreshToken mocks DeleteRefreshToken method
func (m *MockTokenCase) DeleteRefreshToken(ctx context.Context, tokenStr string) error {
	ret := m.Called(ctx, tokenStr)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Error(0)
	}

	return r0
}
