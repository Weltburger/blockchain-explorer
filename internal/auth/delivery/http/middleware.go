package http

import (
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

	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		errJSON := ctx.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("error"))
		return errJSON
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		errJSON := ctx.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("error token format"))
		return errJSON
	}
	if headerParts[0] != "Bearer" {
		errJSON := ctx.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("didn't find bearer"))
		return errJSON
	}

	user, err := m.usecase.ParseToken(ctx.Request().Context(), headerParts[1])
	if err != nil {
		errJSON := ctx.JSON(http.StatusUnauthorized, apperrors.NewAuthorization(err.Error()))
		return errJSON
	}

	ctx.Response().Header().Set(auth.CtxUserKey, user.Email)

	return nil
}

// authorization is the authorization middleware for users.
// It checks the access_token in the Authorization header or the access_token query parameter
// On success sets "me" = *User (current logged user) and "accessData" = current access data
// into the context. Sets even the scopes variable, the sorted slice of scopes in accessData
func Authorization(userUse auth.UserUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return NewAuthMiddleware(userUse)
	}
}
