package main

import (
	"explorer/internal/server"
)

func main() {
	serv := server.NewServer()
	defer serv.Controller.DB.CloseDB()
	//go serv.CheckBlocks()
	go serv.Crawl(690000)
	serv.Router.Logger.Fatal(serv.Router.Start(":1323"))
}
