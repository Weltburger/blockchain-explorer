package http

import (
	"explorer/internal/apperrors"
	"explorer/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
		bindError := err.(*apperrors.Error)
		return c.JSON(bindError.Status(), bindError)
	}

	// validate input fields format and security requirements
	if ok, err := validData(c, req); !ok || err != nil {
		return err
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
		"message": "Account created successfully!",
	})

}
