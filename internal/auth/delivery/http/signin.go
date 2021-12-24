package http

import (
	"explorer/internal/apperrors"
	"explorer/models"
	"log"
	"net/http"

	"github.com/google/uuid"
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
		bindError := err.(*apperrors.Error)
		return c.JSON(bindError.Status(), bindError)
	}

	// validate input fields format and security requirements
	if ok, err := validData(c, req); !ok || err != nil {
		return err
	}

	u := &models.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request().Context()

	err := h.UserUseCase.SignIn(ctx, u)
	if err != nil {
		signError := err.(*apperrors.Error)
		return c.JSON(signError.Status(), signError)
	}

	tokenDetails, err := h.TokenUseCase.NewPairTokens(ctx, u)
	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())
		tokenError := err.(*apperrors.Error)
		return c.JSON(tokenError.Status(), tokenError)
	}

	err = h.TokenUseCase.SavePairTokens(ctx, tokenDetails)
	if err != nil {
		log.Printf("Failed to set tokens to DB: %v\n", err.Error())
		saveErr := err.(*apperrors.Error)
		return c.JSON(saveErr.Status(), saveErr)
	}

	tokens := map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}
	return c.JSON(http.StatusOK, tokens)

}
