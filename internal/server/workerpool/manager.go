package workerpool

import (
	"context"
	"explorer/models"
	"github.com/spf13/viper"
	"time"
)

type Manager struct {
	Counter  int
	Step     int
	TaskChan <-chan int64
	TDataChan chan *TotalData
	ShouldWork bool
	Pool     *Pool
	Data     *TotalData
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
		Step: viper.GetInt("explorer.step"),
		TaskChan: ch,
		TDataChan: make(chan *TotalData),
		ShouldWork: true,
		Pool:     NewPool(nil, viper.GetInt("explorer.totalWorkers")),
		Data: &TotalData{
			Blocks:       make([]models.Block, 0, viper.GetInt("explorer.step")),
			Transactions: make([]models.Transaction, 0),
		},
		Context: ctx,
		Cancel: cancel,
	}
}

func (m *Manager) Process(errChan chan *models.TaskErr,
	dataChan chan *TotalData,
	f func(data int64, dataChan chan *TotalData) error) {
	go m.Pool.RunBackground(m.Context)

	for {
		select {
		case data := <-m.TDataChan:
			m.Data.Blocks = append(m.Data.Blocks, data.Blocks...)
			m.Data.Transactions = append(m.Data.Transactions, data.Transactions...)
			m.Counter++
			if m.Counter == m.Step {
				m.ShouldWork = false
				dataChan <- m.Data
				for !m.ShouldWork {}
			}
		case num := <-m.TaskChan:
			task := NewTask(num, errChan, m.TDataChan, f)
			m.Pool.AddTask(task)
		case <-time.After(5 * time.Second):
			dataChan <- m.Data
			m.Cancel()
		case <-m.Context.Done():
			return
		}
	}
}

func (m *Manager) Reset() {
	m.Counter = 0
	m.Data.Blocks = make([]models.Block, 0, m.Step)
	m.Data.Transactions = make([]models.Transaction, 0)
	m.ShouldWork = true
}


