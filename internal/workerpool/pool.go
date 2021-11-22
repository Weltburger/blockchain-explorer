package workerpool

import (
	"sync"
)

type Pool struct {
	Tasks   []*Task
	Workers []*Worker

	concurrency   int
	Counter       int
	collector     chan *Task
	runBackground chan bool
	Mux			  *sync.Mutex
	WG            *sync.WaitGroup
}

func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:       tasks,
		concurrency: concurrency,
		Counter: 0,
		collector:   make(chan *Task),
		Mux: new(sync.Mutex),
		WG: new(sync.WaitGroup),
	}
}

func (p *Pool) AddTask(task *Task) {
	p.collector <- task
}

func (p *Pool) RunBackground() {
	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.collector, i, p.WG)
		p.Workers = append(p.Workers, worker)
		go worker.StartBackground()
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	p.runBackground = make(chan bool)
	<-p.runBackground
}

func (p *Pool) Stop() {
	for i := range p.Workers {
		p.Workers[i].Stop()
	}
	p.runBackground <- true
}
