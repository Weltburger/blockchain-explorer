package http

import (
	"errors"
	"explorer/internal/auth"
	"explorer/internal/auth/apperrors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	usecase auth.UserUsecase
}

func NewAuthMiddleware(usecase auth.UserUsecase) echo.HandlerFunc {
	return (&AuthMiddleware{
		usecase: usecase,
	}).Handle
}

func (m *AuthMiddleware) Handle(ctx echo.Context) error {

	err := errors.New("middleware")

	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("error"))
		return err
	}

	headerParts := strings.Split(authHeader.(string), " ")
	if len(headerParts) != 2 {
		ctx.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("error token format"))
		return err
	}

	if headerParts[0] != "Bearer" {
		ctx.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("didn't find bearer"))
		return err
	}

	user, err := m.usecase.ParseToken(ctx.Request().Context(), headerParts[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("didn't find bearer"))
		return err
	}

	ctx.Set(auth.CtxUserKey, user)

	return nil
}
