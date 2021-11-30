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
	"time"
)

func Crawl(s *server.Server) {
	wg := &sync.WaitGroup{}
	stepSize := viper.GetInt("explorer.step")
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
					if len(data.Blocks) != stepSize {
						manager.Cancel()
					}
				case err := <-errChan:
					log.Println(err.Err)
					numChan <- err.ID
				case <-ctx.Done():
					break ForLoop
				default:
				}
			}
	}(manager.Context)

	go manager.Process(manager, errChan, dataChan, processingFunc)

	for i := viper.GetInt64("explorer.crawlerStartPos"); i >= 0; i-- {
		numChan <- i
	}

	wg.Wait()
}

func processingFunc(data int64, mng *workerpool.Manager, dataChan chan *workerpool.TotalData) error {
	blockID := data
	for {
		block, err := GetData(&http.Client{}, strconv.FormatInt(blockID, 10))
		if err != nil {
			return err
		}

		if block.Hash != "" {
			transactions := GetTransactions(&block)

			mng.Mux.Lock()
			mng.Counter++
			mng.Data.Blocks = append(mng.Data.Blocks, block)
			mng.Data.Transactions = append(mng.Data.Transactions, transactions...)
			mng.Mux.Unlock()

			if mng.Counter == viper.GetInt("explorer.step") || blockID == 0 {
				time.Sleep(time.Second*10)
				dataChan <- mng.Data
				mng.Reset()
			}

		} else {
			continue
		}

		return nil
	}
}

