package server

import (
	"explorer/internal/explorer/repository/clickhouse"
	"explorer/internal/server/workerpool"
	"explorer/models"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"sync"
	"time"
)

func (s *Server) Crawl(startPos int64) {
	stepSize := viper.GetInt("explorer.step")
	arr := make([]models.Block, 0, stepSize)
	ch := make(chan models.Block)
	chB := make(chan bool)
	chC := make(chan bool)
	chF := make(chan bool)
	chE := make(chan *models.TaskErr)

	defer func() {
		close(ch)
		close(chB)
		close(chC)
		close(chE)
		close(chF)
	}()

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
				ch <- *block

				if pool.Counter == stepSize {
					pool.Counter = 0
					chB <- true
				}
				pool.Mux.Unlock()
			} else {
				continue
			}

			if blockID == 0 {
				chC <- true
			}
			return nil
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		ForLoop:
			for {
				select {
				case str := <-ch:
					arr = append(arr, str)
				case <- chB:
					err := saveData(s, arr...)
					if err != nil {
						log.Println(err)
					}
					arr = make([]models.Block, 0, stepSize)
				case err := <-chE:
					go func() {
						log.Println(err.Err)
						task := workerpool.NewTask(f, err.ID, chE)
						pool.AddTask(task)
					}()
				case <- chC:
					go func() {
						time.Sleep(time.Second * 30)
						err := saveData(s, arr...)
						if err != nil {
							log.Println(err)
						}
						chF <- true
					}()
				case <- chF:
					break ForLoop
				default:
				}
			}
	}()

	go func() {
		defer wg.Done()
		for i := startPos; i >= 0; i-- {
			task := workerpool.NewTask(f, i, chE)
			pool.AddTask(task)
		}
	}()

	wg.Wait()
	pool.Stop()
}

func saveData(s *Server, blocks ...models.Block) error {
	wg := &sync.WaitGroup{}

	br := clickhouse.NewBlockRepository(s.Databases.Clickhouse.DB)
	err := br.PrepareBlockTx()
	if err != nil {
		return err
	}

	defer br.Tx.Rollback()

	wg.Add(2)
	go func(br *clickhouse.BlockRepository) {
		defer wg.Done()
		for _, val := range blocks {
			err := br.Exc(&val)
			if err != nil {
				log.Println(err)
			}
		}
		err := br.Cmt()
		if err != nil {
			log.Println(err)
		}
	}(br)

	go func() {
		defer wg.Done()
		err := saveTransactions(s, blocks...)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	return nil
}

func saveTransactions(s *Server, blocks ...models.Block) error {
	trRepo := clickhouse.NewTransRepository(s.Databases.Clickhouse.DB)
	err := trRepo.PrepareTransactionTx()
	if err != nil {
		return err
	}

	defer trRepo.Tx.Rollback()

	for _, block := range blocks {
		transactions := getTransactions(&block)
		for _, tr := range transactions {
			err := trRepo.SaveTransaction(&tr)
			if err != nil {
				return err
			}
		}
	}

	err = trRepo.Cmt()
	if err != nil {
		return err
	}
	return nil
}

func getTransactions(block *models.Block) []models.Transaction {
	transaction := new(models.Transaction)
	transactions := make([]models.Transaction, 0)
	var opsLength int

	if len(block.Operations) > 2 {
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
	}

	return transactions
}

