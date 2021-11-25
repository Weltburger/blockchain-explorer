package server

import (
	"database/sql"
	"explorer/internal/explorer/delivery/http"
	"explorer/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Router       *echo.Echo
	ClickhouseDB *sql.DB
}

func NewServer() *Server {
	server := &Server{
		Router:     echo.New(),
		ClickhouseDB: storage.GetDB().DB,
	}

	server.Router.Use(middleware.Logger())
	server.Router.Use(middleware.Recover())

	http.RegisterEndpoints(server.Router, server.ClickhouseDB)

	return server
}
