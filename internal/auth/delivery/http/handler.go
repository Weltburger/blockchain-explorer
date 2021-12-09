package http

import (
	"errors"
	"explorer/internal/apperrors"
	"explorer/internal/auth"
	"explorer/models"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	UserUseCase  auth.UserUsecase
	TokenUseCase auth.TokenUsecase
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	UserUsecase  auth.UserUsecase
	TokenUsecase auth.TokenUsecase
}

func NewHandler(c *Config) *Handler {
	return &Handler{
		UserUseCase:  c.UserUsecase,
		TokenUseCase: c.TokenUsecase,
	}
}

type signUpReq struct {
	Email           string `json:"email,omitempty" validate:"required,email"`
	Password        string `json:"password,omitempty" validate:"required,min=8,max=50"`
	ConfirmPassword string `json:"confirm_password,omitempty" validate:"required,min=8,max=50,eqcsfield=Password"`
}

// SignUp handler
func (h *Handler) SignUp(c echo.Context) error {
	// define a variable to which we'll bind incoming json body
	req := new(signUpReq)

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, req); !ok {
		return errors.New("error bind data")
	}

	// validate input fields format and security requirements
	if ok := validData(c, req); !ok {
		return errors.New("error validate data")
	}

	u := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request().Context()

	if err := h.UserUseCase.SignUp(ctx, u); err != nil {
		appErr := err.(*apperrors.Error)
		return c.JSON(apperrors.Status(err), appErr)

	}

	return c.String(http.StatusOK, fmt.Sprintf("Account %s successfully created! Approve your email and Signin!", req.Email))

}

type signInReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// SignIn handler
func (h *Handler) SignIn(c echo.Context) error {
	// define a variable to which we'll bind incoming json body
	req := new(signInReq)

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, req); !ok {
		return errors.New("error bind data")
	}

	// validate input fields format and security requirements
	if ok := validData(c, req); !ok {
		return errors.New("error validate data")
	}

	u := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request().Context()

	err := h.UserUseCase.SignIn(ctx, u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternal())
	}

	tokens, err := h.TokenUseCase.NewPairTokens(ctx, u, "")
	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		return c.JSON(apperrors.Status(err), apperrors.NewAuthorization(err.Error()))
	}

	return c.JSON(http.StatusOK, *tokens)

}
