package main

import (
	"explorer/config"
	"explorer/internal/server"
	"explorer/internal/server/watcher"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

func main() {

	err := config.Init()
	if err != nil {
		log.Printf("error read config: ", err)
		return
	}
	serv, err := server.NewServer()
	if err != nil {
		log.Fatal("error, while initialising server: ", err)
		return
	}

	defer serv.Databases.Clickhouse.DB.Close()
	defer serv.Databases.Postgres.DB.Close()

	//go watcher.CheckBlocks(serv)
	go watcher.Crawl(serv)

	serv.Router.Logger.Fatal(serv.Router.Start(viper.GetString("address")))
}
