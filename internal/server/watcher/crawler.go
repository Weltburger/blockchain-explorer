package watcher

import (
	"context"
	"explorer/internal/server"
	"explorer/internal/server/workerpool"
	"explorer/models"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func Crawl(s *server.Server) {
	wg := &sync.WaitGroup{}
	//stepSize := viper.GetInt("explorer.step")
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
					fmt.Println("aboba")
					//saveBlocks(s, data.Blocks...)
					//saveTransactions(s, data.Transactions...)
					fmt.Println(len(data.Blocks))
					fmt.Println(len(data.Transactions))
					manager.Reset()
					fmt.Println(len(data.Blocks))
					fmt.Println(len(data.Transactions))
					/*if len(data.Blocks) != stepSize {
						manager.Cancel()
					}*/

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
				/*if blockID == 0 {
					time.Sleep(time.Second*5)
				}*/
				dataChan <- td
			}()

			/*mng.Mux.Lock()
			mng.Counter++
			mng.Data.Blocks = append(mng.Data.Blocks, block)
			mng.Data.Transactions = append(mng.Data.Transactions, transactions...)


			if mng.Counter == viper.GetInt("explorer.step") || blockID == 0 {
				time.Sleep(time.Second*5)
				dataChan <- mng.Data
				time.Sleep(time.Second)
				mng.Reset()
			}
			mng.Mux.Unlock()*/

		} else {
			continue
		}

		return nil
	}
	return nil
}
