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
	"strconv"
	"sync"
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

	apiGroup.GET("/transactions", server.Controller.TransactionController().GetTransactions)
	apiGroup.GET("/transactions/block/:block", server.Controller.TransactionController().GetTransactionsByBlock)
	apiGroup.GET("/transactions/address/:address", server.Controller.TransactionController().GetTransactionsByAddress)
	apiGroup.GET("/transactions/hash/:hash", server.Controller.TransactionController().GetTransactionsByHash)

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

		err = s.Controller.DB.BlockStorage().SaveBlock(&block)
		if err != nil {
			log.Fatal(err)
		}
		err = s.parseTransactions(&block)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(block)

		time.Sleep(time.Second * 30)
	}
}

func (s *Server) Crawl(startPos uint64) {
	wg := &sync.WaitGroup{}

	for startPos > 1 {
		url := "https://mainnet-tezos.giganode.io/chains/main/blocks/" + strconv.FormatUint(startPos, 10)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		block, _ := models.UnmarshalBlock(body)

		wg.Add(2)
		go func() {
			defer wg.Done()
			err = s.Controller.DB.BlockStorage().SaveBlock(&block)
			if err != nil {
				log.Fatal(err)
			}
		}()
		go func() {
			defer wg.Done()
			err = s.parseTransactions(&block)
			if err != nil {
				log.Fatal(err)
			}
		}()
		wg.Wait()

		fmt.Println(block)
		startPos--
		time.Sleep(time.Second * 1)
	}
}

func (s *Server) parseTransactions(block *models.Block) error {
	transaction := new(models.Transaction)

	for i := 0; i < len(block.Operations[3]); i++ {
		for j := 0; j < len(block.Operations[3][i].Contents); j++ {
			transaction.BlockHash = block.Hash
			transaction.Hash = block.Operations[3][i].Hash
			transaction.Branch = block.Operations[3][i].Branch
			transaction.Destination = block.Operations[3][i].Contents[j].Destination
			transaction.Source = block.Operations[3][i].Contents[j].Source
			transaction.Fee = block.Operations[3][i].Contents[j].Fee
			transaction.Counter = block.Operations[3][i].Contents[j].Counter
			transaction.GasLimit = block.Operations[3][i].Contents[j].GasLimit
			transaction.Amount = block.Operations[3][i].Contents[j].Amount
			transaction.ConsumedMilligas = block.Operations[3][i].Contents[j].Metadata.OperationResult.ConsumedMilligas
			transaction.StorageSize = block.Operations[3][i].Contents[j].Metadata.OperationResult.StorageSize
			transaction.Signature = block.Operations[3][i].Signature

			err := s.Controller.DB.TransactionStorage().SaveTransaction(transaction)
			if err != nil {
				return err
			}
		}
	}
	err := s.Controller.DB.TransactionStorage().Cmt()
	if err != nil {
		return err
	}
	return nil
}
