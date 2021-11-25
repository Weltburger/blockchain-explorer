package server

import (
	authhttp "explorer/internal/auth/delivery/http"
	authrepo "explorer/internal/auth/repository/postgres"
	"explorer/internal/auth/usecase"
	"explorer/internal/controller"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
)

func inject(d *dataSources) *Server {
	log.Println("Injecting data sources...")

	// repository layer
	userRepo := authrepo.NewUserRepository(d.DB)

	// usecase(service) layer
	userUseCase := usecase.NewAuthUseCase(userRepo, []byte("password"), 600)

	// initialize Echo router
	router := echo.New()

	// use default middleware
	router.Use(middleware.Logger())
	router.Logger.SetLevel(echolog.DEBUG)
	router.Use(middleware.Recover())

	// create input fields validator
	validate := validator.New()

	// register auth endpoints
	authhttp.RegisterEndpoints(router, userUseCase, validate)
	// authMiddleware := authhttp.NewAuthMiddleware(server.AuthUC)

	// create server struct
	server := &Server{
		Router:     router,
		Controller: controller.New(),
		AuthUC:     userUseCase,
	}

	return server
}
