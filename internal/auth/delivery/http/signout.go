package http

import (
	"explorer/internal/apperrors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary SignOut
// @Tags auth
// @Description Logout from the service
// @Produce  json
// @ID sign-out
// @Success 200 {object} map[string]string
// @Failure 400 {object} apperrors.Error
// @Failure 401 {object} apperrors.Error
// @Failure 500 {object} apperrors.Error
// @Router /api/auth/sign-out [post]
func (h *Handler) SignOut(c echo.Context) error {

	ctx := c.Request().Context()
	requestHeader := c.Request()

	err := h.TokenUseCase.DeleteTokens(ctx, requestHeader)
	if err != nil {
		appErr := err.(*apperrors.Error)
		return c.JSON(appErr.Status(), appErr)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logout successfully!",
	})
}
