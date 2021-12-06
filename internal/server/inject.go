package server

import (
	authhttp "explorer/internal/auth/delivery/http"
	pgrepo "explorer/internal/auth/repository/postgres"
	redisrepo "explorer/internal/auth/repository/redis"
	"explorer/internal/auth/usecase"
	explhttp "explorer/internal/explorer/delivery/http"
	"explorer/internal/storage"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
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
	userUC := usecase.NewUserUseCase(userRepo)

	// load rsa keys
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := ioutil.ReadFile(privKeyFile)

	if err != nil {
		return nil, fmt.Errorf("Could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		return nil, fmt.Errorf("Could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)

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

	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse ID_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	tokenUC := usecase.NewTokenUsecase(&usecase.TSConfig{
		TokenRepository:       tokenRepo,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         refreshSecret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	// tokenTTL, err := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	// if err != nil {
	// 	return nil, err
	// }

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

	api := server.Router.Group("/api")

	// register auth endpoints
	authhttp.RegisterEndpoints(api, authhttp.Config{UserUsecase: server.UserUC, TokenUsecase: server.TokenUC})
	api.Use(authhttp.Authorization(server.TokenUC))

	// register explorer endpoints
	explhttp.RegisterEndpoints(server.Router, ds.Clickhouse.DB)

	return server, nil
}
