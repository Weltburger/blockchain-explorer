package http

import (
	"github.com/labstack/echo/v4"
)

func RegisterEndpoints(router *echo.Group, conf Config) {
	h := NewHandler(&conf)

	// authEndpoints := router.Group("/auth")
	{
		router.POST("/sign-up", h.SignUp)
		router.POST("/sign-in", h.SignIn)
	}

}
