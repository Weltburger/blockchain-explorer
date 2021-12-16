package server

import (
	"explorer/internal/auth"
	"explorer/internal/storage"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Router    *echo.Echo
	UserUC    auth.UserUsecase
	TokenUC   auth.TokenUsecase
	Databases *storage.DataSources
}

func NewServer() (*Server, error) {
	server, err := inject()
	if err != nil {
		return nil, err
	}

	return server, nil
}
