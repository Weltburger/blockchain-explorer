package workerpool

import (
	"explorer/models"
	"fmt"
)

type Task struct {
	Err  error
	Data int64
	f    func(int64) error
	chE  chan *models.TaskErr
}

func NewTask(f func(int64) error, data int64, ch chan *models.TaskErr) *Task {
	return &Task{f: f, Data: data, chE: ch}
}

func process(workerID int, task *Task) {
	fmt.Printf("Worker %d processes task %v\n", workerID, task.Data)
	task.Err = task.f(task.Data)
	if task.Err != nil {
		task.chE <- &models.TaskErr{
			Err: task.Err,
			ID:  task.Data,
		}
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	}
}
