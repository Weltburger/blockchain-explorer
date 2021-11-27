package main

import (
	"explorer/internal/server"
)

func main() {
	serv := server.NewServer()
	defer serv.ClickhouseDB.Close()
	go serv.CheckBlocks()
	go serv.Crawl(678500, 500)
	serv.Router.Logger.Fatal(serv.Router.Start(":1323"))
}
