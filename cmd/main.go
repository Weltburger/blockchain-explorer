package main

import (
	"explorer/internal/server"
)

func main() {
	serv := server.NewServer()
	defer serv.Controller.DB.CloseDB()
	//go serv.CheckBlocks()
	go serv.Crawl(250, 200/*678500*/)
	serv.Router.Logger.Fatal(serv.Router.Start(":1323"))
}
