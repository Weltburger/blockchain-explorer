package http

import (
	"explorer/internal/auth"

	"github.com/labstack/echo/v4"
)

func RegisterEndpoints(router *echo.Echo, uc auth.UserUsecase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.POST("/sign-in", h.SignIn)
	}
}
