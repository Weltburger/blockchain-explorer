package workerpool

import (
	"context"
	"fmt"
)

type Worker struct {
	ID       int
	taskChan chan *Task
	context      context.Context
}

func NewWorker(channel chan *Task, ID int, ctx context.Context) *Worker {
	return &Worker{
		ID:       ID,
		taskChan: channel,
		context:  ctx,
	}
}

func (wr *Worker) StartBackground() {
	fmt.Printf("Starting worker %d\n", wr.ID)

	for {
		select {
		case task := <-wr.taskChan:
			process(wr.ID, task)
		case <-wr.context.Done():
			fmt.Printf("Closing worker %d\n", wr.ID)
			return
		}
	}
}
