package http

import (
	"explorer/internal/apperrors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type refreshReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// @Summary Refresh
// @Tags auth
// @Description refresh token pair
// @Accept json
// @Produce  json
// @ID refresh
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
		return err
	}

	ctx := c.Request().Context()

	refreshDetails, err := h.TokenUseCase.ValidateToken(ctx, req.RefreshToken, true)
	if err != nil {
		validRefError := err.(*apperrors.Error)
		return c.JSON(validRefError.Status(), validRefError)
	}

	newTokenPair, err := h.TokenUseCase.GenerateTokens(ctx, refreshDetails.UserId)
	if err != nil {
		log.Printf("Error generate new token pair: %v", err)
		return c.JSON(apperrors.NewInternal().Status(), apperrors.NewInternal())
	}

	err = h.TokenUseCase.DeleteRefreshToken(ctx, refreshDetails.TokenId)
	if err != nil {
		delError := err.(*apperrors.Error)
		return c.JSON(delError.Status(), delError)
	}

	return c.JSON(http.StatusOK, newTokenPair)

}
