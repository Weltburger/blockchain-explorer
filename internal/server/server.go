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
	Router     *echo.Echo
	Controller *controller.Controller
}

func NewServer() *Server {
	server := &Server{
		Router:     echo.New(),
		Controller: controller.New(),
	}

	server.Router.Use(middleware.Logger())
	server.Router.Use(middleware.Recover())

	apiGroup := server.Router.Group("/api/v1")

	apiGroup.GET("/blocks", server.Controller.BlockController().GetBlocks)
	apiGroup.GET("/block/:block", server.Controller.BlockController().GetBlock)

	apiGroup.GET("/transactions", server.Controller.TransactionController().GetTransactions)


	return server
}

func (s *Server) CheckBlocks() {
	for {
		block, err := getData("head")
		if err != nil {
			log.Fatal(err)
		}

		//saveData(s, block)

		fmt.Println(*block)
		time.Sleep(time.Second * 30)
	}
}

func (s *Server) Crawl(startPos int64) {
	var step int64 = 10000
	arr := make([]*models.Block, 0, step)
	ch := make(chan *models.Block, 1000)
	chB := make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case str := <-ch:
				arr = append(arr, str)
			case b := <- chB:
				if b == true {
					fmt.Println(startPos)
					startPos -= step
					for _, val := range arr{
						saveData(s, val)
					}
					if startPos < 10000 {
						arr = make([]*models.Block, 0, startPos)
						go crawlStep(startPos, &ch, &chB)
					} else {
						fmt.Println("Done", len(arr), arr)
						arr = make([]*models.Block, 0, step)
						go crawlLess(startPos, &ch, &chB)
					}
				} else {
					return
				}
			default:
				break
			}
		}
	}()

	if startPos < step {
		go crawlLess(startPos, &ch, &chB)
	} else {
		go crawlStep(startPos, &ch, &chB)
	}
	wg.Wait()

	for _, val := range arr{
		saveData(s, val)
	}
}

func crawlStep(start int64, ch *chan *models.Block, chB *chan bool) {
	wg := new(sync.WaitGroup)

	for i := 1; i < 11; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			for j := 1; j < 11; {

				block, err := getData(strconv.FormatInt(start-int64(i*j), 10))
				if err != nil {
					log.Fatal(err)
				}
				if block.Hash != "" {
					fmt.Println(*block)
					j++
					*ch <- block
				}
			}
		}(wg, i)
	}

	wg.Wait()
	time.Sleep(time.Millisecond*50)
	*chB <- true
}

func crawlLess(start int64, ch *chan *models.Block, chB *chan bool) {
	wg := new(sync.WaitGroup)
	for i := 1; i < int(start)+1; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()

			block, err := getData(strconv.FormatInt(start-int64(i), 10))
			if err != nil {
				log.Fatal(err)
			}

			*ch <- block
		}(wg, i)
	}

	wg.Wait()
	time.Sleep(time.Millisecond*50)
	*chB <- false
}

func getData(index string) (*models.Block, error) {
	url := "https://testnet-tezos.giganode.io/chains/main/blocks/" + index
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	block, _ := models.UnmarshalBlock(body)

	return &block, nil
}

func saveData(s *Server, block *models.Block) {
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		err := s.Controller.DB.BlockStorage().SaveBlock(block)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := s.parseTransactions(block)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

func (s *Server) parseTransactions(block *models.Block) error {
	transactionStorage := s.Controller.DB.TransactionStorage()
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

			err := transactionStorage.SaveTransaction(transaction)
			if err != nil {
				return err
			}
		}
	}
	err := transactionStorage.Cmt()
	if err != nil {
		return err
	}
	return nil
}
