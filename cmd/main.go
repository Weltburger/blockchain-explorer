package main

import (
	"explorer/internal/server"
	"os"

	"github.com/labstack/gommon/log"
)

func main() {

	// err := config.Init()
	// if err != nil {
	// 	log.Printf("error read config: ", err)
	// 	return
	// }
	serv, err := server.NewServer()
	if err != nil {
		log.Fatal("error, while initializing server: ", err)
		return
	}

	defer serv.Databases.Clickhouse.DB.Close()
	defer serv.Databases.Postgres.DB.Close()

	go serv.CheckBlocks()
	go serv.Crawl(678500)

	serv.Router.Logger.Fatal(serv.Router.Start(os.Getenv("HTTP_PORT")))
}
