package server

import (
	"explorer/internal/controller"
	"explorer/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Router *echo.Echo
	Controller *controller.Controller
}

func NewServer() *Server {
	server := &Server{
		Router: echo.New(),
		Controller: controller.New(),
	}

	server.Router.Use(middleware.Logger())
	server.Router.Use(middleware.Recover())

	apiGroup := server.Router.Group("/api/v1")

	apiGroup.GET("/blocks", server.Controller.BlockController().GetBlocks)
	apiGroup.GET("/block/:block", server.Controller.BlockController().GetBlock)

	return server
}

func (s *Server) CheckBlocks() {
	for {
		resp, err := http.Get("https://mainnet-tezos.giganode.io/chains/main/blocks/head")
		if err != nil {
			log.Fatal(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		block, _ := models.UnmarshalBlock(body)

		/*err = s.Controller.DB.BlockStorage().SaveBlock(&block)
		if err != nil {
			log.Fatal(err)
		}*/
		//fmt.Println(block)
		for i := 0; i < len(block.Operations[3]); i++ {
			fmt.Println(block.Operations[3][i])
		}

		time.Sleep(time.Second * 30)
	}
}

func (s *Server) Crawl(startPosHash string) {
	for {
		// правильное распараллеливание + ограничение по запросам?
		url := "https://mainnet-tezos.giganode.io/chains/main/blocks/" + startPosHash + "-"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		block, _ := models.UnmarshalBlock(body)

		/*err = s.Controller.DB.BlockStorage().SaveBlock(&block)
		if err != nil {
			log.Fatal(err)
		}*/
		fmt.Println(block)

		time.Sleep(time.Second * 30)
	}
}
