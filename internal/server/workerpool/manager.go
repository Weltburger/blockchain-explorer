package workerpool

import (
	"context"
	"explorer/models"
	"log"
	"os"
	"strconv"
)

type Manager struct {
	Counter    int64
	Total      int64
	Step       int64
	TaskChan   <-chan models.CrawlRange
	TDataChan  chan *TotalData
	ErrChan    chan *models.TaskErr
	ErrHistory map[int64]int
	ShouldWork bool
	OuterQueue bool
	Pool       *Pool
	Data       *TotalData
	Context    context.Context
	Cancel     context.CancelFunc
}

type TotalData struct {
	Blocks         []models.Block
	Transactions   []models.Transaction
	TransactionsMI []models.TransactionMainInfo
}

func CreateManager(ch <-chan models.CrawlRange) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	step, err := strconv.ParseInt(os.Getenv("STEP"), 10, 64)
	if err != nil {
		log.Fatal("error, while getting crawler step: ", err)
	}

	totalWorkers, err := strconv.Atoi(os.Getenv("TOTAL_WORKERS"))
	if err != nil {
		log.Fatal("error, while getting manager total workers: ", err)
	}
	return &Manager{
		Counter:    0,
		Total:      0,
		Step:       step,
		TaskChan:   ch,
		TDataChan:  make(chan *TotalData),
		ErrChan:    make(chan *models.TaskErr),
		ErrHistory: map[int64]int{},
		ShouldWork: true,
		OuterQueue: false,
		Pool:       NewPool(nil, totalWorkers),
		Data: &TotalData{
			Blocks:       make([]models.Block, 0, step),
			Transactions: make([]models.Transaction, 0),
			TransactionsMI: make([]models.TransactionMainInfo, 0),
		},
		Context: ctx,
		Cancel:  cancel,
	}
}

func (m *Manager) Process(dataChan chan *TotalData,
	f func(data int64, dataChan chan *TotalData) error) {
	go m.Pool.RunBackground(m.Context)

	end, err := strconv.ParseInt(os.Getenv("CRAWLER_END_POS"), 10, 64)
	if err != nil {
		log.Fatal("error, while getting crawler ending position: ", err)
	}
	start, err := strconv.ParseInt(os.Getenv("CRAWLER_START_POS"), 10, 64)
	if err != nil {
		log.Fatal("error, while getting crawler starting position: ", err)
	}
	totalOps := (end - start) + 1

	for {
		select {
		case data := <-m.TDataChan:
			m.Data.Blocks = append(m.Data.Blocks, data.Blocks...)
			m.Data.Transactions = append(m.Data.Transactions, data.Transactions...)
			m.Data.TransactionsMI = append(m.Data.TransactionsMI, data.TransactionsMI...)
			m.Counter++
			m.Total++
			if m.Counter == m.Step || m.Total == totalOps {
				m.ShouldWork = false
				dataChan <- m.Data
				for !m.ShouldWork {
				}
			}
		case num := <-m.TaskChan:
			go func() {
				for i := num.From; i < num.To; i++ {
					task := NewTask(i, m.ErrChan, m.TDataChan, f)
					m.Pool.AddTask(task)
				}
			}()
		case err := <-m.ErrChan:
			log.Println(err.Err)
			if res, ok := m.ErrHistory[err.ID]; ok {
				if res == 4 {
					log.Println("error occurred five times, so we'll skip it")
					m.Counter++
					totalOps--
					continue
				}
				m.ErrHistory[err.ID]++
			} else {
				m.ErrHistory[err.ID] = 1
			}
			task := NewTask(err.ID, m.ErrChan, m.TDataChan, f)
			m.Pool.AddTask(task)
		case <-m.Context.Done():
			return
		}
	}
}

func (m *Manager) Reset() {
	m.Counter = 0
	m.Data.Blocks = make([]models.Block, 0, m.Step)
	m.Data.Transactions = make([]models.Transaction, 0)
	m.Data.TransactionsMI = make([]models.TransactionMainInfo, 0)
	m.ShouldWork = true
	m.OuterQueue = false
}
