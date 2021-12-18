package main

import (
	"explorer/internal/server"
	"explorer/internal/server/watcher"
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
)

// @title Blockchain Explorer
// @version 0.9.1
// @description This is a service that allows you to receive data stored in the blockchain.

// @host localhost
// @BasePath /api

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	serv, err := server.NewServer()
	if err != nil {
		log.Fatal("error, while initializing server: ", err)
		return
	}

	defer serv.Databases.Clickhouse.DB.Close()
	defer serv.Databases.Postgres.DB.Close()
	defer serv.Databases.Redis.Close()

	go watcher.CheckBlocks(serv)
	//go watcher.Crawl(serv)

	serv.Router.Logger.Fatal(serv.Router.Start(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))))
}
