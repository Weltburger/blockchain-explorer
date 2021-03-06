package http

import (
	"explorer/internal/apperrors"
	"explorer/internal/auth"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	token auth.TokenUsecase
	next  echo.HandlerFunc
}

func NewAuthMiddleware(t auth.TokenUsecase, next echo.HandlerFunc) echo.HandlerFunc {
	return (&AuthMiddleware{
		token: t,
		next:  next,
	}).Handle
}

func (m *AuthMiddleware) Handle(c echo.Context) error {

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("Header authorization field is empty."))
	}

	tokenParts := strings.Split(authHeader, "Bearer ")
	if len(tokenParts) != 2 {
		return c.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`"))
	}
	if strings.Contains(tokenParts[0], "Bearer") {
		return c.JSON(http.StatusUnauthorized, apperrors.NewAuthorization("Didn't find 'Bearer'"))
	}

	ctx := c.Request().Context()

	user, err := m.token.ValidateToken(ctx, tokenParts[1], false)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, apperrors.NewAuthorization(err.Error()))
	}

	c.Set(auth.CtxUserKey, user)

	return m.next(c)
}

// authorization is the authorization middleware for users.
// It checks the access_token in the Authorization header or the access_token query parameter
func Authorization(t auth.TokenUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return NewAuthMiddleware(t, next)
	}
}
