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

	//apiGroup.GET("/blocks", Controller.UserController().CreateUser)
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

		/*s.Controller.DB.BlockStorage().Exc(&block)
		s.Controller.DB.BlockStorage().Cmt()*/
		fmt.Println(block)

		time.Sleep(time.Second * 15)
	}
}
