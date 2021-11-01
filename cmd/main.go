package main

import (
	"explorer/internal/server"
)

func main() {
	serv := server.NewServer()
	defer serv.Controller.DB.CloseDB()
	go serv.CheckBlocks()
	serv.Router.Logger.Fatal(serv.Router.Start(":1323"))
}
