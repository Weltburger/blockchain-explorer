package workerpool

import (
	"context"
	"explorer/models"
	"github.com/spf13/viper"
	"sync"
)

type Manager struct {
	Counter  int
	TaskChan <-chan int64
	Pool     *Pool
	Data     *TotalData
	Mux      *sync.Mutex
	Context  context.Context
	Cancel   context.CancelFunc
}

type TotalData struct {
	Blocks       []models.Block
	Transactions []models.Transaction
}

func CreateManager(ch <-chan int64) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		Counter:  0,
		TaskChan: ch,
		Pool:     NewPool(nil, viper.GetInt("explorer.totalWorkers")),
		Data: &TotalData{
			Blocks:       make([]models.Block, 0, viper.GetInt("explorer.step")),
			Transactions: make([]models.Transaction, 0),
		},
		Mux: new(sync.Mutex),
		Context: ctx,
		Cancel: cancel,
	}
}

func (m *Manager) Process(mng *Manager, errChan chan *models.TaskErr, dataChan chan *TotalData, f func(data int64, mng *Manager, dataChan chan *TotalData) error) {
	go m.Pool.RunBackground(m.Context)

	for {
		select {
		case num := <-m.TaskChan:
			task := NewTask(mng, num, errChan, dataChan, f)
			m.Pool.AddTask(task)
		case <-m.Context.Done():
			return
		}
	}
}

func (m *Manager) Reset() {
	m.Mux.Lock()
	m.Counter = 0
	m.Data.Blocks = make([]models.Block, 0, viper.GetInt("explorer.step"))
	m.Data.Transactions = make([]models.Transaction, 0)
	m.Mux.Unlock()
}


