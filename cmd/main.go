package main

import (
	"explorer/internal/server"
	"explorer/internal/server/watcher"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
)

func main() {
	serv, err := server.NewServer()
	if err != nil {
		log.Fatal("error, while initializing server: ", err)
		return
	}

	defer serv.Databases.Clickhouse.DB.Close()
	defer serv.Databases.Postgres.DB.Close()

	go watcher.CheckBlocks(serv)
	go watcher.Crawl(serv)

	serv.Router.Logger.Fatal(serv.Router.Start(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))))
}
