package http

import (
	"explorer/internal/apperrors"
	"explorer/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type refreshReq struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// @Summary Refresh
// @Tags auth
// @Description refresh token pair
// @Accept json
// @Produce  json
// @ID refresh
// @Param access_token body string true "access token"
// @Param refresh_token body string true "refresh token"
// @Success 200 {object} models.TokenPair
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Failure 500 {object} apperrors.Error
// @Router /api/auth/refresh [post]
func (h *Handler) Refresh(c echo.Context) error {

	// define a variable to which we'll bind incoming json body
	req := new(refreshReq)

	// Bind incoming json to struct and check for validation errors
	if ok, err := bindData(c, req); !ok || err != nil {
		bindError := err.(*apperrors.Error)
		return c.JSON(bindError.Status(), bindError)
	}

	// validate input fields format and security requirements
	if ok, err := validData(c, req); !ok || err != nil {
		validError := err.(*apperrors.Error)
		return c.JSON(validError.Status(), validError)
	}

	ctx := c.Request().Context()

	refreshDetails, err := h.TokenUseCase.ValidateRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		validRefError := err.(*apperrors.Error)
		return c.JSON(validRefError.Status(), validRefError)
	}

	newTokenDetails, err := h.TokenUseCase.NewPairTokens(ctx, refreshDetails.UserId)
	if err != nil {
		log.Printf("Error generate new token pair: %v", err)
		return c.JSON(apperrors.NewInternal().Status(), apperrors.NewInternal())
	}

	err = h.TokenUseCase.SavePairTokens(ctx, newTokenDetails)
	if err != nil {
		log.Printf("Failed to set tokens to DB: %v\n", err.Error())
		saveErr := err.(*apperrors.Error)
		return c.JSON(saveErr.Status(), saveErr)
	}

	err = h.TokenUseCase.DeleteRefreshToken(ctx, refreshDetails.TokenUuid)
	if err != nil {
		delError := err.(*apperrors.Error)
		return c.JSON(delError.Status(), delError)
	}

	tokens := &models.TokenPair{
		AccessToken:  newTokenDetails.AccessToken,
		RefreshToken: newTokenDetails.RefreshToken,
	}

	return c.JSON(http.StatusOK, tokens)

}
