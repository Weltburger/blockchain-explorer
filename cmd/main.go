package main

import (
	"explorer/internal/server"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"strconv"
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

	startPos, err := strconv.ParseInt(os.Getenv("CRAWLER_START_POS"), 10, 64)
	if err != nil {
		log.Fatal("error, while getting crawler start position: ", err)
		return
	}
	go serv.CheckBlocks()
	go serv.Crawl(startPos)

	serv.Router.Logger.Fatal(serv.Router.Start(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))))
}
