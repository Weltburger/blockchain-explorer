package watcher

import (
	"context"
	"explorer/internal/server"
	"explorer/internal/server/workerpool"
	"explorer/models"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func Crawl(s *server.Server) {
	wg := &sync.WaitGroup{}
	dataChan := make(chan *workerpool.TotalData)
	errChan := make(chan *models.TaskErr)
	numChan := make(chan int64)

	defer func() {
		close(dataChan)
		close(errChan)
		close(numChan)
	}()

	manager := workerpool.CreateManager(numChan)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		ForLoop:
			for {
				select {
				case data := <-dataChan:
					saveBlocks(s, data.Blocks...)
					saveTransactions(s, data.Transactions...)
					manager.Reset()
				case err := <-errChan:
					log.Println(err.Err)
					numChan <- err.ID
				case <-ctx.Done():
					break ForLoop
				default:
				}
			}
	}(manager.Context)

	go manager.Process(errChan, dataChan, processingFunc)

	for i := viper.GetInt64("explorer.crawlerStartPos"); i >= 0; i-- {
		numChan <- i
	}

	wg.Wait()
}

func processingFunc(data int64, dataChan chan *workerpool.TotalData) error {
	blockID := data
	for {
		block, err := GetData(&http.Client{}, strconv.FormatInt(blockID, 10))
		if err != nil {
			return err
		}

		if block.Hash != "" {
			transactions := GetTransactions(&block)

			td := &workerpool.TotalData{
				Blocks:       []models.Block{block},
				Transactions: transactions,
			}

			go func() {
				dataChan <- td
			}()
		} else {
			continue
		}

		return nil
	}
}
