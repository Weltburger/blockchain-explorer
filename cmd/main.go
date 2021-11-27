package main

import (
	"explorer/config"
	"explorer/internal/server"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

func main() {

	err := config.Init()
	if err != nil {
		log.Printf("error read config: %v", err)
		return
	}
	serv := server.NewServer()

	defer serv.ClickhouseDB.Close()
	go serv.CheckBlocks()
	go serv.Crawl(678500, 500)
	serv.Router.Logger.Fatal(serv.Router.Start(viper.GetString("address")))
}
