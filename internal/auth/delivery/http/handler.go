package http

import (
	"explorer/internal/apperrors"
	"explorer/internal/auth"
	"explorer/models"
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

// @Summary SignUp
// @Tags auth
// @Description Register user in the service
// @Accept  json
// @Produce  json
// @ID sign-up
// @Param email body string true "user email/login"
// @Param password body string true "user password"
// @Param confirm_password body string true "confirm user password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} apperrors.Error
// @Failure 409 {object} apperrors.Error
// @Failure 500 {object} apperrors.Error
// @Router /api/auth/sign-up [post]
func (h *Handler) SignUp(c echo.Context) error {
	// define a variable to which we'll bind incoming json body
	req := new(signUpReq)

	// Bind incoming json to struct and check for validation errors
	if ok, err := bindData(c, req); !ok || err != nil {
		bindError := apperrors.NewBadRequest(err.Error())
		return c.JSON(bindError.Status(), bindError)
	}

	// validate input fields format and security requirements
	if ok, err := validData(c, req); !ok || err != nil {
		validError := apperrors.NewBadRequest(err.Error())
		return c.JSON(validError.Status(), validError)
	}

	u := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request().Context()

	if err := h.UserUseCase.SignUp(ctx, u); err != nil {
		appErr := err.(*apperrors.Error)
		return c.JSON(appErr.Status(), appErr)

	}

	return c.JSON(http.StatusOK, map[string]string{
		"type":    "SUCCESSES",
		"message": "Account successfully created!",
	})

}

type signInReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// @Summary SignIn
// @Tags auth
// @Description Authorize user in the service
// @Accept  json
// @Produce  json
// @ID sign-in
// @Param email body string true "user login"
// @Param password body string true "user password"
// @Success 200 {object} models.TokenPair
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Failure 404 {object} apperrors.Error
// @Failure 409 {object} apperrors.Error
// @Failure 500 {object} apperrors.Error
// @Router /api/auth/sign-ip [post]
func (h *Handler) SignIn(c echo.Context) error {
	// define a variable to which we'll bind incoming json body
	req := new(signInReq)

	// Bind incoming json to struct and check for validation errors
	if ok, err := bindData(c, req); !ok || err != nil {
		bindError := apperrors.NewBadRequest(err.Error())
		return c.JSON(bindError.Status(), bindError)
	}

	// validate input fields format and security requirements
	if ok, err := validData(c, req); !ok || err != nil {
		validError := apperrors.NewBadRequest(err.Error())
		return c.JSON(validError.Status(), validError)
	}

	u := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request().Context()

	err := h.UserUseCase.SignIn(ctx, u)
	if err != nil {
		signError := err.(*apperrors.Error)
		return c.JSON(signError.Status(), signError)
	}

	tokens, err := h.TokenUseCase.NewPairTokens(ctx, u, "")
	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())
		tokenError := err.(*apperrors.Error)
		return c.JSON(tokenError.Status(), tokenError)
	}

	return c.JSON(http.StatusOK, *tokens)

}
