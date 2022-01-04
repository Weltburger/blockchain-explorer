package http

import (
	"explorer/internal/apperrors"
	"explorer/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
// @Param email body string true "user email/login"
// @Param password body string true "user password"
// @Success 200 {object} models.TokenPair
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Failure 404 {object} apperrors.Error
// @Failure 409 {object} apperrors.Error
// @Failure 500 {object} apperrors.Error
// @Router /api/auth/sign-in [post]
func (h *Handler) SignIn(c echo.Context) error {
	// define a variable to which we'll bind incoming json body
	req := new(signInReq)

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

	err := h.UserUseCase.SignIn(ctx, u)
	if err != nil {
		signError := err.(*apperrors.Error)
		return c.JSON(signError.Status(), signError)
	}

	tokens, err := h.TokenUseCase.GenerateTokens(ctx, u.ID.String())
	if err != nil {
		log.Printf("Failed to create tokens for userID: %s - %v\n", u.ID.String(), err)
		tokenError := err.(*apperrors.Error)
		return c.JSON(tokenError.Status(), tokenError)
	}

	return c.JSON(http.StatusOK, tokens)

}
