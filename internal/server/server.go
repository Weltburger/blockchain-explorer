package server

import (
	"explorer/internal/auth"
	authhttp "explorer/internal/auth/delivery/http"
	authrepo "explorer/internal/auth/repository/postgres"
	"explorer/internal/auth/usecase"
	"explorer/internal/controller"
	"explorer/internal/storage"
	"explorer/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Router     *echo.Echo
	Controller *controller.Controller
	AuthUC     auth.UserUsecase
}

func NewServer() *Server {
	ds, err := initDS()
	if err != nil {
		log.Println("error init DB connect")
		return nil
	}
	userRepo := authrepo.NewUserRepository(ds.DB)
	userUseCase := usecase.NewAuthUseCase(userRepo, "salt", []byte("pass"), time.Second*100)

	server := &Server{
		Router:     echo.New(),
		Controller: controller.New(),
		AuthUC:     userUseCase,
	}

	server.Router.Use(middleware.Logger())
	server.Router.Use(middleware.Recover())

	authhttp.RegisterEndpoints(server.Router, server.AuthUC)
	authMiddleware := authhttp.NewAuthMiddleware(server.AuthUC)
	apiGroup := server.Router.Group("/api/v1", authMiddleware)

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
	mux := &sync.Mutex{}
	var step int64 = 10000
	arr := make([]*models.Block, 0, step)
	ch := make(chan *models.Block /*, 1000*/)
	chB := make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(mux *sync.Mutex) {
		defer wg.Done()
		for {
			select {
			case str := <-ch:
				mux.Lock()
				arr = append(arr, str)
				mux.Unlock()
			case b := <-chB:
				if b == true {
					fmt.Println(startPos)
					startPos -= step
					saveDataRange(s, arr)
					if startPos < 10000 {
						arr = make([]*models.Block, 0, startPos)
						go crawlLess(startPos, &ch, &chB)
					} else {
						arr = make([]*models.Block, 0, step)
						go crawlStep(startPos, &ch, &chB)
					}
				} else {
					return
				}
			default:
				break
			}
		}
	}(mux)

	if startPos < step {
		go crawlLess(startPos, &ch, &chB)
	} else {
		go crawlStep(startPos, &ch, &chB)
	}
	wg.Wait()

	saveDataRange(s, arr)
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
	//time.Sleep(time.Millisecond*50)
	*chB <- true
}

func crawlLess(start int64, ch *chan *models.Block, chB *chan bool) {
	wg := new(sync.WaitGroup)
	for i := 1; i < int(start)+1; {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()

			block, err := getData(strconv.FormatInt(start-int64(i), 10))
			if err != nil {
				log.Fatal(err)
			}

			if block.Hash != "" {
				fmt.Println(*block)
				i++
				*ch <- block
			}
		}(wg, i)
	}

	wg.Wait()
	//time.Sleep(time.Millisecond*50)
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
		err := saveTransactions(s, block)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

func saveDataRange(s *Server, blocks []*models.Block) {
	wg := &sync.WaitGroup{}
	bs := s.Controller.DB.BlockStorage()

	wg.Add(2)
	go func(bs *storage.BlockStorage) {
		defer wg.Done()
		for _, val := range blocks {
			err := bs.Exc(val)
			if err != nil {
				log.Fatal(err)
			}
		}
		err := bs.Cmt()
		if err != nil {
			log.Fatal(err)
		}
	}(bs)

	go func() {
		defer wg.Done()
		err := saveRangeTransactions(s, blocks)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

func saveTransactions(s *Server, block *models.Block) error {
	transactionStorage := s.Controller.DB.TransactionStorage()

	transactions := getTransactions(block)

	for _, tr := range transactions {
		err := transactionStorage.Exc(tr)
		if err != nil {
			return err
		}
	}
	err := transactionStorage.Cmt()
	if err != nil {
		return err
	}
	return nil
}

func saveRangeTransactions(s *Server, blocks []*models.Block) error {
	transactionStorage := s.Controller.DB.TransactionStorage()

	for _, block := range blocks {
		transactions := getTransactions(block)
		for _, tr := range transactions {
			err := transactionStorage.Exc(tr)
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

func getTransactions(block *models.Block) []*models.Transaction {
	transaction := new(models.Transaction)
	transactions := make([]*models.Transaction, 0)
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

			transactions = append(transactions, transaction)
		}
	}

	return transactions
}
