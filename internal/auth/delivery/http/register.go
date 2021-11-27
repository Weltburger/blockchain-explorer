package http

import (
	"explorer/internal/auth"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func RegisterEndpoints(router *echo.Group, uc auth.UserUsecase, v *validator.Validate) {
	h := NewHandler(uc, v)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.POST("/sign-in", h.SignIn)
	}

}
