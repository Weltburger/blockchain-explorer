package http

import "context"

// @Summary SignOut
// @Tags auth
// @Description Logout from the service
// @Produce  json
// @ID sign-out
// @Success 200 {object} string
// @Failure 400 {object} apperrors.Error
// @Failure 500 {object} apperrors.Error
// @Router /api/v1/sign-out [post]
func (h *Handler) SignOut(ctx context.Context) error {
	return nil
}
