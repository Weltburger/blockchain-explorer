package workerpool

import (
	"context"
	"sync"
)

type Pool struct {
	Tasks   []*Task
	Workers []*Worker

	concurrency   int
	Counter       int
	collector     chan *Task
	Mux			  *sync.Mutex
}

func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:       tasks,
		concurrency: concurrency,
		Counter: 0,
		collector:   make(chan *Task),
		Mux: new(sync.Mutex),
	}
}

func (p *Pool) AddTask(task *Task) {
	p.collector <- task
}

func (p *Pool) RunBackground(ctx context.Context) {
	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.collector, i, ctx)
		p.Workers = append(p.Workers, worker)
		go worker.StartBackground()
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}
}
