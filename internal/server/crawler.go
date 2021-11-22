package server

import (
	"explorer/internal/storage"
	"explorer/internal/workerpool"
	"explorer/models"
	"fmt"
	"log"
	"strconv"
	"sync"
)

func (s *Server) Crawl(startPos int64) {
	arr := make([]models.Block, 0, 50)
	fail := make([]int64, 0)
	ch := make(chan models.Block)
	chB := make(chan bool)
	chE := make(chan *models.TaskErr)

	pool := workerpool.NewPool(nil, 10)
	go pool.RunBackground()

	f := func(data int64) error {
		blockID := data
		for {
			block, err := getData(strconv.FormatInt(blockID, 10))
			if err != nil {
				return err
			}

			if block.Hash != "" {
				pool.Mux.Lock()
				pool.Counter++
				fmt.Println(*block)
				ch <- *block

				if pool.Counter == 50 {
					pool.Counter = 0
					chB <- true
				}
				pool.Mux.Unlock()
			} else {
				continue
			}

			if blockID == 1 {
				chB <- false
			}
			return nil
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			select {
			case str := <-ch:
				arr = append(arr, str)
			case b := <- chB:
				if b == true {
					err := saveData(s, arr...)
					if err != nil {
						return
					}
				}
			case err := <-chE:
				fmt.Println(err)
				fail = append(fail, err.ID)
				//go finish(left, ff, pool, chE, chB)
				//fail = make([]int64, 0)
			default:
				//break
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := startPos; i > 0; i-- {
			task := workerpool.NewTask(f, i, chE)
			pool.AddTask(task)
		}
	}()

	wg.Wait()
	pool.Stop()
}

func crawlStep(start int64, ch chan models.Block, chB chan bool) {
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
					ch <- *block
				}
			}
		}(wg, i)
	}

	wg.Wait()
	//time.Sleep(time.Millisecond*50)
	chB <- true
}

func crawlLess(start int64, ch chan models.Block, chB chan bool) {
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
				ch <- *block
			}
		}(wg, i)
	}

	wg.Wait()
	//time.Sleep(time.Millisecond*50)
	chB <- false
}

func saveData(s *Server, blocks ...models.Block) error {
	wg := &sync.WaitGroup{}
	bs := s.Controller.DB.BlockStorage()
	err := bs.PrepareBlockTx()
	if err != nil {
		return err
	}

	wg.Add(2)
	go func(bs *storage.BlockStorage) {
		defer wg.Done()
		for _, val := range blocks {
			err := bs.Exc(&val)
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
		err := saveTransactions(s, blocks...)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	return nil
}

func saveTransactions(s *Server, blocks ...models.Block) error {
	transactionStorage := s.Controller.DB.TransactionStorage()
	err := transactionStorage.PrepareTransactionTx()
	if err != nil {
		return err
	}


	for _, block := range blocks {
		transactions := getTransactions(&block)
		for _, tr := range transactions {
			err := transactionStorage.SaveTransaction(&tr)
			if err != nil {
				return err
			}
		}
	}

	err = transactionStorage.Cmt()
	if err != nil {
		return err
	}
	return nil
}

func getTransactions(block *models.Block) []models.Transaction {
	transaction := new(models.Transaction)
	transactions := make([]models.Transaction, 0)
	var opsLength int

	if opsLength = len(block.Operations[3]); opsLength == 0 {
		return nil
	}

	for i := 0; i < opsLength; i++ {
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

			transactions = append(transactions, *transaction)
		}
	}

	return transactions
}

