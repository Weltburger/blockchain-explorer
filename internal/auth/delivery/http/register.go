package http

import (
	"github.com/labstack/echo/v4"
)

func RegisterEndpoints(router *echo.Group, conf Config) {
	h := NewHandler(&conf)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.POST("/sign-in", h.SignIn)
		authEndpoints.POST("/refresh", h.Refresh)
		authEndpoints.POST("/sign-out", h.SignOut, Authorization(h.TokenUseCase))
	}

}
