package watcher

import (
	"explorer/internal/explorer/repository/clickhouse"
	"explorer/internal/server"
	"explorer/models"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Getter interface {
	Get(url string) (*http.Response, error)
}

func GetData(getter Getter, index string, prevPos int64) (models.Block, error) {
	var (
		resp *http.Response
		err error
		block models.Block
	)
	for {
		url := fmt.Sprintf(fmt.Sprintf("%s%s", os.Getenv("NODE"), index))
		resp, err = getter.Get(url)
		if err != nil {
			return models.Block{}, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return models.Block{}, err
		}
		block, _ = models.UnmarshalBlock(body)

		if block.Hash == "" {
			continue
		} else if block.Metadata.LevelInfo.Level == prevPos {
			time.Sleep(time.Second * 5)
			continue
		} else {
			return block, nil
		}
	}
}

func saveAllData(s *server.Server, blocks ...models.Block) error {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := saveBlocks(s, blocks...)
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := saveTransactionsFromBlocks(s, blocks...)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	return nil
}

func saveBlocks(s *server.Server, blocks ...models.Block) error {
	blRepo := clickhouse.NewBlockRepository(s.Databases.Clickhouse.DB)
	err := blRepo.PrepareBlockTx()
	if err != nil {
		return err
	}

	defer blRepo.Rollback()

	for _, block := range blocks {
		err := blRepo.Exc(&block)
		if err != nil {
			return err
		}
	}
	err = blRepo.Cmt()
	if err != nil {
		return err
	}

	return nil
}

func saveTransactions(s *server.Server, transactions ...models.Transaction) error {
	trRepo := clickhouse.NewTransRepository(s.Databases.Clickhouse.DB)
	err := trRepo.PrepareTransactionTx()
	if err != nil {
		return err
	}

	defer trRepo.Tx.Rollback()

	for _, transaction := range transactions {
		err := trRepo.Exc(&transaction)
		if err != nil {
			return err
		}
	}
	err = trRepo.Cmt()
	if err != nil {
		return err
	}

	return nil
}

func saveTransactionsMI(s *server.Server, transactions ...models.TransactionMainInfo) error {
	trRepo := clickhouse.NewTransMIRepository(s.Databases.Clickhouse.DB)
	err := trRepo.PrepareTransactionTx()
	if err != nil {
		return err
	}

	defer trRepo.Tx.Rollback()

	for _, transaction := range transactions {
		err := trRepo.Exc(&transaction)
		if err != nil {
			return err
		}
	}
	err = trRepo.Cmt()
	if err != nil {
		return err
	}

	return nil
}

func saveTransactionsFromBlocks(s *server.Server, blocks ...models.Block) error {
	trRepo := clickhouse.NewTransRepository(s.Databases.Clickhouse.DB)
	trMIRepo := clickhouse.NewTransMIRepository(s.Databases.Clickhouse.DB)

	err := trRepo.PrepareTransactionTx()
	if err != nil {
		return err
	}

	err = trMIRepo.PrepareTransactionTx()
	if err != nil {
		return err
	}

	defer trRepo.Tx.Rollback()

	for _, block := range blocks {
		transactions, transactionsMI := GetTransactions(&block)
		for _, tr := range transactions {
			err := trRepo.SaveTransaction(&tr)
			if err != nil {
				return err
			}
		}
		for _, tr := range transactionsMI {
			err := trMIRepo.SaveTransaction(&tr)
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

func GetTransactions(block *models.Block) ([]models.Transaction, []models.TransactionMainInfo) {
	transaction := models.Transaction{}
	transactions := make([]models.Transaction, 0)
	transactionMI := models.TransactionMainInfo{}
	transactionsMI := make([]models.TransactionMainInfo, 0)
	var opsLength int

	if len(block.Operations) > 2 {
		if opsLength = len(block.Operations[3]); opsLength == 0 {
			return nil, nil
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
				transactions = append(transactions, transaction)

				if transaction.Source == viper.GetString("explorer.contract") ||
					transaction.Destination == viper.GetString("explorer.contract") {
					transactionMI.BlockHash = transaction.BlockHash
					transactionMI.Hash = transaction.Hash
					transactionMI.Source = transaction.Source
					transactionMI.Destination = transaction.Destination
					transactionMI.Fee = transaction.Fee
					transactionMI.Amount = transaction.Amount
					transactionsMI = append(transactionsMI, transactionMI)
				}
			}
		}
	}

	return transactions, transactionsMI
}

