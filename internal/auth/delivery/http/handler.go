package http

import (
	"explorer/internal/auth"
)

type Handler struct {
	UserUseCase  auth.UserUsecase
	TokenUseCase auth.TokenUsecase
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	UserUsecase  auth.UserUsecase
	TokenUsecase auth.TokenUsecase
}

func NewHandler(c *Config) *Handler {
	return &Handler{
		UserUseCase:  c.UserUsecase,
		TokenUseCase: c.TokenUsecase,
	}
}
