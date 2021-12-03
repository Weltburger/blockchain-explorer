package server

import (
	authhttp "explorer/internal/auth/delivery/http"
	authrepo "explorer/internal/auth/repository/postgres"
	"explorer/internal/auth/usecase"
	explhttp "explorer/internal/explorer/delivery/http"
	"explorer/internal/storage"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
)

func inject() (*Server, error) {
	log.Println("Injecting data sources...")

	DataSources, err := storage.InitDataSources()
	if err != nil {
		return nil, err
	}

	tokenTTL, err := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	server := &Server{
		Router: echo.New(),
		AuthUC: usecase.NewAuthUseCase(authrepo.NewUserRepository(DataSources.Postgres.DB),
			[]byte(os.Getenv("SIGNING_KEY")), time.Duration(tokenTTL)),
		Databases: DataSources,
	}

	server.Router.Debug = true

	// use default middleware
	server.Router.Use(middleware.Logger())
	server.Router.Logger.SetLevel(echolog.DEBUG)
	server.Router.Use(middleware.Recover())

	server.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// create input fields validator
	validate := validator.New()

	api := server.Router.Group("/api")

	// register explorer endpoints
	explhttp.RegisterEndpoints(server.Router, DataSources.Clickhouse.DB)

	// register auth endpoints
	authhttp.RegisterEndpoints(api, server.AuthUC, validate)
	api.Use(authhttp.Authorization(server.AuthUC))

	return server, nil
}
