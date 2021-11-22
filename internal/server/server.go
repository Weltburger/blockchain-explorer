package server

import (
	"explorer/internal/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Router     *echo.Echo
	Controller *controller.Controller
}

func NewServer() *Server {
	server := &Server{
		Router:     echo.New(),
		Controller: controller.New(),
	}

	server.Router.Use(middleware.Logger())
	server.Router.Use(middleware.Recover())

	apiGroup := server.Router.Group("/api/v1")

	apiGroup.GET("/blocks", server.Controller.BlockController().GetBlocks)
	apiGroup.GET("/block/:block", server.Controller.BlockController().GetBlock)

	apiGroup.GET("/transactions", server.Controller.TransactionController().GetTransactions)


	return server
}
