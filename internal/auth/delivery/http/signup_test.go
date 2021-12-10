package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"explorer/internal/apperrors"
	"explorer/models"
	"explorer/models/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	e := echo.New()

	mockUC := new(mocks.MockUserCase)
	mockTC := new(mocks.MockTokenCase)

	// create router group
	r := e.Group("/api")

	RegisterEndpoints(r, Config{UserUsecase: mockUC, TokenUsecase: mockTC})

	t.Run("Email and Password required", func(t *testing.T) {
		signUpBody := signUpReq{
			Email: "1@mail.local",
		}

		mockUC.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*models.User")).Return(nil)

		body, err := json.Marshal(signUpBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		e.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUC.AssertNotCalled(t, "SignUp", mock.Anything)
	})

	t.Run("Invalid Email", func(t *testing.T) {
		signUpBody := signUpReq{
			Email:           "1@mail",
			Password:        "password",
			ConfirmPassword: "password",
		}

		mockUC.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*models.User")).Return(nil)

		body, err := json.Marshal(signUpBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		e.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockUC.AssertNotCalled(t, "SignUp", mock.Anything)
	})

	t.Run("Password too shot", func(t *testing.T) {
		signUpBody := signUpReq{
			Email:           "1@mail.local",
			Password:        "123",
			ConfirmPassword: "123",
		}

		mockUC.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*models.User")).Return(nil)

		body, err := json.Marshal(signUpBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		e.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockUC.AssertNotCalled(t, "SignUp", mock.Anything)
	})

	t.Run("Invalid header Content-Type", func(t *testing.T) {
		signUpBody := signUpReq{
			Email:           "1@mail.local",
			Password:        "1qaz2wsx3edc",
			ConfirmPassword: "1qaz2wsx3edc",
		}

		mockUC.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*models.User")).Return(nil)

		body, err := json.Marshal(signUpBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "plain/text")

		e.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)

		mockUC.AssertNotCalled(t, "SignUp", mock.Anything)
	})

	t.Run("Invalid JSON body format", func(t *testing.T) {
		wrongJSONBody := `{"email": "1@mail.local", "password": "password"`

		mockUC.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*models.User")).Return(nil)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBuffer([]byte(wrongJSONBody)))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		e.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUC.AssertNotCalled(t, "SignUp", mock.Anything)
	})

	t.Run("Error returned from UserUsecase", func(t *testing.T) {
		e := echo.New()

		mockUC := new(mocks.MockUserCase)
		mockTC := new(mocks.MockTokenCase)

		// create router group
		r := e.Group("/api")

		RegisterEndpoints(r, Config{UserUsecase: mockUC, TokenUsecase: mockTC})

		u := &models.User{
			Email:    "1@mail.local",
			Password: "strongpassword",
		}

		mockUC.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), u).Return(apperrors.NewConflict("email", u.Email))

		signUpBody := signUpReq{
			Email:           u.Email,
			Password:        u.Password,
			ConfirmPassword: u.Password,
		}
		body, err := json.Marshal(signUpBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		e.ServeHTTP(w, req)

		mockConflictError := apperrors.NewConflict("email", u.Email)
		// respBody, err := json.Marshal(mockConflictError)
		// assert.NoError(t, err)

		assert.Equal(t, mockConflictError.Status(), w.Code)
		// assert.Equal(t, respBody, w.Body.Bytes())
		mockUC.AssertExpectations(t)
	})

	t.Run("Successfull sign-up", func(t *testing.T) {
		e := echo.New()

		mockUC := new(mocks.MockUserCase)
		mockTC := new(mocks.MockTokenCase)

		// create router group
		r := e.Group("/api")

		RegisterEndpoints(r, Config{UserUsecase: mockUC, TokenUsecase: mockTC})

		u := &models.User{
			Email:    "1@mail.local",
			Password: "strongpassword",
		}

		mockUC.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), u).Return(nil)

		signUpBody := signUpReq{
			Email:           u.Email,
			Password:        u.Password,
			ConfirmPassword: u.Password,
		}
		body, err := json.Marshal(signUpBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/auth/sign-up", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		e.ServeHTTP(w, req)

		responseBody := fmt.Sprintf("Account %s successfully created! Approve your email and Signin!", u.Email)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, responseBody, w.Body.String())
		mockUC.AssertExpectations(t)
	})

}
