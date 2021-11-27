package server

import (
	"database/sql"
	"explorer/internal/auth"
	"github.com/labstack/echo/v4"
	"log"
)

type Server struct {
	Router     *echo.Echo
	AuthUC     auth.UserUsecase
	ClickhouseDB *sql.DB
}

func NewServer() *Server {
	ds, err := initDS()
	if err != nil {
		log.Println("error init DB connect")
		return nil
	}

	server := inject(ds)

	return server
}