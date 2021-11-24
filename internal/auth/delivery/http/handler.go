package http

import (
	"explorer/internal/auth"
	"explorer/internal/auth/apperrors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase auth.UserUsecase
	// valid   *validator.Validate
}

func NewHandler(useCase auth.UserUsecase) *Handler {
	return &Handler{
		useCase: useCase,
		// valid:   v,
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
		ctx.JSON(http.StatusBadRequest, apperrors.NewBadRequest(err.Error()))
		return err
	}

	if err := h.useCase.SignUp(ctx.Request().Context(), input.Email, input.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, apperrors.NewInternal())
		return err
	}

	ctx.JSON(http.StatusOK, "OK")
	return nil
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
		ctx.JSON(http.StatusBadRequest, apperrors.NewBadRequest(err.Error()))
		return err
	}

	token, err := h.useCase.SignIn(ctx.Request().Context(), input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apperrors.NewInternal())
		return err
	}

	ctx.JSON(http.StatusOK, signInResponse{Token: token})

	return nil
}
