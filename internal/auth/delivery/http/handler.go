package http

import (
	"explorer/internal/apperrors"
	"explorer/internal/auth"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase  auth.UserUsecase
	validate *validator.Validate
}

func NewHandler(useCase auth.UserUsecase, v *validator.Validate) *Handler {
	return &Handler{
		useCase:  useCase,
		validate: v,
	}
}

type signUpData struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=50"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=50,eqcsfield=Password"`
}

func (h *Handler) SignUp(ctx echo.Context) error {
	input := new(signUpData)

	if err := ctx.Bind(input); err != nil {
		errJSON := ctx.JSON(http.StatusBadRequest, apperrors.NewBadRequest(err.Error()))
		return errJSON
	}

	// validate input data
	if err := h.validate.Struct(input); err != nil {
		errValidator := make([]string, 0, 3)
		for _, e := range err.(validator.ValidationErrors) {
			errValidator = append(errValidator, strings.Split(e.Error(), "Error:")[1])
		}
		errJSON := ctx.JSON(http.StatusBadRequest, apperrors.NewBadRequest(strings.Join(errValidator, " --> ")))
		return errJSON
	}

	if err := h.useCase.SignUp(ctx.Request().Context(), input.Email, input.Password); err != nil {
		errJSON := ctx.JSON(http.StatusInternalServerError, apperrors.NewInternal())
		return errJSON
	}

	errJSON := ctx.JSON(http.StatusOK, "OK")
	return errJSON
}

type signInData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

type signInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(ctx echo.Context) error {
	input := new(signInData)

	if err := ctx.Bind(input); err != nil {
		errJSON := ctx.JSON(http.StatusBadRequest, apperrors.NewBadRequest(err.Error()))
		return errJSON
	}

	// validate input data
	if err := h.validate.Struct(input); err != nil {
		errValidator := make([]string, 0, 2)
		for _, e := range err.(validator.ValidationErrors) {
			errValidator = append(errValidator, strings.Split(e.Error(), "Error:")[1])
		}
		errJSON := ctx.JSON(http.StatusBadRequest, apperrors.NewBadRequest(strings.Join(errValidator, " --> ")))
		return errJSON
	}

	token, err := h.useCase.SignIn(ctx.Request().Context(), input.Email, input.Password)
	if err != nil {
		errJSON := ctx.JSON(http.StatusInternalServerError, apperrors.NewInternal())
		return errJSON
	}

	errJSON := ctx.JSON(http.StatusOK, signInResponse{Token: token})
	return errJSON
}
