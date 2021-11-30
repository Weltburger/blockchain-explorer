package workerpool

import (
	"explorer/models"
	"fmt"
)

type Task struct {
	Err  error
	Data int64
	ErrChan  chan *models.TaskErr
	DataChan chan *TotalData
	Manager  *Manager
	processingFunc func(data int64, mng *Manager, dataChan chan *TotalData) error
}

func NewTask(mng *Manager,
	data int64,
	errChan chan *models.TaskErr,
	dataChan chan *TotalData,
	f func(data int64, mng *Manager, dataChan chan *TotalData) error) *Task {
	return &Task{Data: data,
		ErrChan: errChan,
		DataChan: dataChan,
		Manager: mng,
		processingFunc: f,
	}
}

func process(workerID int, task *Task) {
	fmt.Printf("Worker %d processes task %v\n", workerID, task.Data)
	task.Err = task.processingFunc(task.Data, task.Manager, task.DataChan)
	if task.Err != nil {
		task.ErrChan <- &models.TaskErr{
			Err: task.Err,
			ID:  task.Data,
		}
	} else {
		fmt.Printf("Worker %d completed task %v\n", workerID, task.Data)
	}
}
