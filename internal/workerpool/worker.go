package workerpool

import (
	"fmt"
	"sync"
)

type Worker struct {
	ID       int
	taskChan chan *Task
	quit     chan bool
	WG       *sync.WaitGroup
}

func NewWorker(channel chan *Task, ID int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:       ID,
		taskChan: channel,
		quit:     make(chan bool),
		WG: wg,
	}
}

func (wr *Worker) StartBackground() {
	fmt.Printf("Starting worker %d\n", wr.ID)

	for {
		select {
		case task := <-wr.taskChan:
			wr.WG.Add(1)
			process(wr.ID, task, wr.WG)
		case <-wr.quit:
			return
		}
	}
}

func (wr *Worker) Stop() {
	fmt.Printf("Closing worker %d\n", wr.ID)
	go func() {
		wr.quit <- true
	}()
}
