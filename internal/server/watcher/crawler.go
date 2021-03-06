package watcher

import (
	"context"
	"explorer/internal/server"
	"explorer/internal/server/workerpool"
	"explorer/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func Crawl(s *server.Server) {
	wg := &sync.WaitGroup{}
	dataChan := make(chan *workerpool.TotalData)
	rangeChan := make(chan models.CrawlRange)

	defer func() {
		close(dataChan)
		close(rangeChan)
	}()

	manager := workerpool.CreateManager(rangeChan)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case data := <-dataChan:
				saveBlocks(s, data.Blocks...)
				saveTransactions(s, data.Transactions...)
				saveTransactionsMI(s, data.TransactionsMI...)
				manager.Reset()
			case <-ctx.Done():
				return
			default:
			}
		}
	}(manager.Context)

	go manager.Process(dataChan, processingFunc)

	start, err := strconv.ParseInt(os.Getenv("CRAWLER_START_POS"), 10, 64)
	if err != nil {
		log.Fatal("error, while getting crawler start position: ", err)
		return
	}
	end, err := strconv.ParseInt(os.Getenv("CRAWLER_END_POS"), 10, 64)
	if err != nil {
		log.Fatal("error, while getting crawler start position: ", err)
		return
	}
	step, err := strconv.ParseInt(os.Getenv("STEP"), 10, 64)
	if err != nil {
		log.Fatal("error, while getting crawler start position: ", err)
		return
	}

	if start > end {
		start, end = end, start
	}
	cRange := models.CrawlRange{}
	for i := start; i <= end; {
		if manager.OuterQueue {
			continue
		}
		cRange.From = i
		if (i+step) > end {
			cRange.To = end + 1
		} else {
			cRange.To = i + step
		}
		rangeChan <- cRange

		i+=step
		manager.OuterQueue = true
	}

	wg.Wait()
}

func processingFunc(data int64, dataChan chan *workerpool.TotalData) error {
	blockID := data
	for {
		block, err := GetData(&http.Client{}, strconv.FormatInt(blockID, 10), -1)
		if err != nil {
			return err
		}

		if block.Metadata.LevelInfo.Level == 0 {
			block.Metadata.LevelInfo = block.Metadata.Level
		}

		transactions, transactionsMI := GetTransactions(&block)

		td := &workerpool.TotalData{
			Blocks:       []models.Block{block},
			Transactions: transactions,
			TransactionsMI: transactionsMI,
		}

		go func() {
			dataChan <- td
		}()

		return nil
	}
}
