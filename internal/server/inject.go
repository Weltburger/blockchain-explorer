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
	"github.com/spf13/viper"
)

func inject(d *dataSources) *Server {
	log.Println("Injecting data sources...")

	// repository layer
	userRepo := authrepo.NewUserRepository(d.DB)

	// usecase(service) layer
	userUseCase := usecase.NewAuthUseCase(userRepo, []byte(viper.GetString("auth.signing_key")), viper.GetDuration("auth.token_ttl"))

	// initialize Echo router
	router := echo.New()
	router.Debug = true

	// use default middleware
	router.Use(middleware.Logger())
	router.Logger.SetLevel(echolog.DEBUG)
	router.Use(middleware.Recover())

	// create input fields validator
	validate := validator.New()

	api := router.Group("/api")

	// register auth endpoints
	authhttp.RegisterEndpoints(api, userUseCase, validate)
	api.Use(authhttp.Authorization(userUseCase))

	// api.GET("/test", func(c echo.Context) error {
	// 	return c.JSON(http.StatusOK, "Hello, World!")
	// })

	// create server struct
	server := &Server{
		Router:     router,
		Controller: controller.New(),
		AuthUC:     userUseCase,
	}

	return server
}
