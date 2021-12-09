package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"explorer/internal/apperrors"
	"explorer/internal/auth/usecase/mocks"
	"explorer/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignIn(t *testing.T) {
	e := echo.New()

	mockUC := new(mocks.MockUserCase)
	mockTC := new(mocks.MockTokenCase)

	// create router group
	r := e.Group("/api")

	RegisterEndpoints(r, Config{UserUsecase: mockUC, TokenUsecase: mockTC})

	t.Run("Invalid request data", func(t *testing.T) {
		signUpBody := signInReq{
			Email:    "1@mail",
			Password: "strongpassword",
		}

		body, err := json.Marshal(signUpBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		e.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockUC.AssertNotCalled(t, "SignIn", mock.Anything)
		mockTC.AssertNotCalled(t, "NewPairTokens", mock.Anything)
	})

	t.Run("Invalid request data JSON format", func(t *testing.T) {
		wrongJSONBody := `{"email": "1@mail.local", "password": "password"`

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBuffer([]byte(wrongJSONBody)))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		e.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertNotCalled(t, "SignIn", mock.Anything)
		mockTC.AssertNotCalled(t, "NewPairTokens", mock.Anything)
	})

	t.Run("Error returned from UserUsecase SignIn", func(t *testing.T) {
		signInBody := signInReq{
			Email:    "1@mail.local",
			Password: "strongpassword",
		}

		body, err := json.Marshal(signInBody)
		assert.NoError(t, err)

		u := &models.User{
			Email:    signInBody.Email,
			Password: signInBody.Password,
		}
		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			u,
		}

		// so we can check for a known status code
		mockNotFoundError := apperrors.NewNotFound("email", u.Email)
		// mockNotAuthError := apperrors.NewAuthorization("Invalid email and password combination")
		mockUC.On("SignIn", mockUSArgs...).Return(mockNotFoundError)
		// mockUC.On("SignIn", mockUSArgs...).Return(mockNotAuthError)

		// a response recorder for getting written http response
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		e.ServeHTTP(w, req)

		mockUC.AssertCalled(t, "SignIn", mockUSArgs...)
		mockTC.AssertNotCalled(t, "NewPairTokens")
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		e := echo.New()

		mockUC := new(mocks.MockUserCase)
		mockTC := new(mocks.MockTokenCase)

		// create router group
		r := e.Group("/api")

		RegisterEndpoints(r, Config{UserUsecase: mockUC, TokenUsecase: mockTC})

		signInReq := signInReq{
			Email:    "1@mail.local",
			Password: "strongpassword",
		}

		u := &models.User{
			Email:    signInReq.Email,
			Password: signInReq.Password,
		}
		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			u,
		}
		mockUC.On("SignIn", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			u,
			"",
		}
		mockTokenPair := &models.TokenPair{
			IDToken:      models.IDToken{SS: "idToken"},
			RefreshToken: models.RefreshToken{SS: "refreshToken"},
		}
		mockTC.On("NewPairTokens", mockTSArgs...).Return(mockTokenPair, nil)

		// a response recorder for getting written http response
		w := httptest.NewRecorder()

		// create a request body with valid fields
		body, err := json.Marshal(signInReq)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		e.ServeHTTP(w, req)

		// respBody, err := json.Marshal(mockTokenPair)
		// assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		// assert.Equal(t, respBody, w.Body.Bytes())

		mockUC.AssertCalled(t, "SignIn", mockUSArgs...)
		mockTC.AssertCalled(t, "NewPairTokens", mockTSArgs...)
	})

	t.Run("Error Token Creation", func(t *testing.T) {
		e := echo.New()

		mockUC := new(mocks.MockUserCase)
		mockTC := new(mocks.MockTokenCase)

		// create router group
		r := e.Group("/api")

		RegisterEndpoints(r, Config{UserUsecase: mockUC, TokenUsecase: mockTC})

		signInReq := signInReq{
			Email:    "1@mail.local",
			Password: "strongpassword",
		}

		u := &models.User{
			Email:    signInReq.Email,
			Password: signInReq.Password,
		}
		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			u,
		}
		mockUC.On("SignIn", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			u,
			"",
		}

		mockError := apperrors.NewInternal()
		mockTC.On("NewPairTokens", mockTSArgs...).Return(nil, mockError)

		// a response recorder for getting written http response
		w := httptest.NewRecorder()

		// create a request body with valid fields
		body, err := json.Marshal(signInReq)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-in", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		e.ServeHTTP(w, req)

		// respBody, err := json.Marshal(mockError)
		// assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), w.Code)
		// assert.Equal(t, respBody, w.Body.Bytes())

		mockUC.AssertCalled(t, "SignIn", mockUSArgs...)
		mockTC.AssertCalled(t, "NewPairTokens", mockTSArgs...)
	})

}
