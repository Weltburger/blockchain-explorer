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
	defer serv.Controller.DB.CloseDB()
	//go serv.CheckBlocks()
	// go serv.Crawl(9000)
	serv.Router.Logger.Fatal(serv.Router.Start(viper.GetString("address")))
}
