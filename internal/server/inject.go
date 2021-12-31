package server

import (
	authhttp "explorer/internal/auth/delivery/http"
	pgrepo "explorer/internal/auth/repository/postgres"
	redisrepo "explorer/internal/auth/repository/redis"
	"explorer/internal/auth/usecase"
	explhttp "explorer/internal/explorer/delivery/http"
	"explorer/internal/storage"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"

	_ "explorer/docs"
)

func inject() (*Server, error) {
	log.Println("Injecting data sources...")

	ds, err := storage.InitDataSources()
	if err != nil {
		return nil, err
	}

	/*
	 * repository layer
	 */
	userRepo := pgrepo.NewUserRepository(ds.Postgres.DB)
	tokenRepo := redisrepo.NewTokenRepository(ds.Redis.Client)

	/*
	 * service layer
	 */
	userUC := usecase.NewUserCase(userRepo)

	// load rsa keys
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := os.ReadFile(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		return nil, fmt.Errorf("Could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := os.ReadFile(pubKeyFile)

	if err != nil {
		return nil, fmt.Errorf("Could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	if err != nil {
		return nil, fmt.Errorf("Could not parse public key: %w", err)
	}

	// load refresh token secret from env variable
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// load expiration lengths from env variables and parse as int
	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	idExp, err := strconv.ParseInt(idTokenExp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse ID_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	tokenUC := usecase.NewTokenUsecase(&usecase.TSConfig{
		TokenRepository:       tokenRepo,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         refreshSecret,
		AccessExpirationSecs:  idExp,
		RefreshExpirationSecs: refreshExp,
	})

	server := &Server{
		Router:    echo.New(),
		UserUC:    userUC,
		TokenUC:   tokenUC,
		Databases: ds,
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

	server.Router.GET("/swagger/*", echoSwagger.WrapHandler)

	api := server.Router.Group("/api")

	// register auth endpoints
	authhttp.RegisterEndpoints(api, authhttp.Config{UserUsecase: server.UserUC, TokenUsecase: server.TokenUC})

	// set middleware
	api.Use(authhttp.Authorization(server.TokenUC))

	// register explorer endpoints
	explhttp.RegisterEndpoints(api, ds.Clickhouse.DB)

	return server, nil
}
