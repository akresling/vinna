package vinna

import (
	"fmt"
	"time"
)

// Pool is the interface for the worker pool manager
type Pool interface {
	StartAllWorkers()
	StopAllWorkers()
	Stop()
	AddWorker(int)
	DropWorker()
	ManageWorkers()
}

// WorkerPool is a worker manager without responses
type WorkerPool struct {
	initNumWorker int
	nWorkers      int
	Jobs          <-chan Work
	workers       []Worker
	running       bool
	sigStop       chan bool
}

// NewWorkerPool returns a new worker pool manager
func NewWorkerPool(n int, jobs <-chan Work) WorkerPool {
	return WorkerPool{
		initNumWorker: n,
		nWorkers:      0,
		Jobs:          jobs,
		workers:       make([]Worker, 0),
		running:       false,
		sigStop:       make(chan bool),
	}
}

// StartAllWorkers will initialize all the workers in the pool
func (p *WorkerPool) StartAllWorkers() {
	if p.running {
		return
	}
	for i := 1; i <= p.initNumWorker; i++ {
		p.AddWorker(i)
	}
	p.running = true
	go p.ManageWorkers()
}

// StopAllWorkers will stop all the workers in the pool
func (p *WorkerPool) StopAllWorkers() {
	for _, g := range p.workers {
		g.Stop()
	}
}

// Stop signals the manager to stop
func (p *WorkerPool) Stop() {
	p.sigStop <- true
	p.StopAllWorkers()
	p.running = false
}

// AddWorker will add a new worker to the pool and start it
func (p *WorkerPool) AddWorker(id int) {
	w := NewGenericWorker(id, p.Jobs)
	p.workers = append(p.workers, &w)
	p.nWorkers++
	w.Start()
}

// DropWorker will stop a worker  and remove it from the pool
func (p *WorkerPool) DropWorker() {
	if p.nWorkers == 1 {
		return
	}
	w := p.workers[len(p.workers)-1]
	w.Stop()
	p.workers = p.workers[:len(p.workers)-1]
	p.nWorkers--
}

// ManageWorkers will control the number of workers in the pool
func (p *WorkerPool) ManageWorkers() {
	for {
		select {
		case <-time.Tick(50 * time.Millisecond):
			if len(p.Jobs) > cap(p.Jobs)/5 {
				p.AddWorker(p.nWorkers + 1)
			} else if len(p.Jobs) < cap(p.Jobs)/10 {
				p.DropWorker()
			}
		case <-p.sigStop:
			fmt.Println("STOPPING MANAGER")
			return
		}
	}
}

func debug(p interface{}) {
	fmt.Printf("%v\n", p)
}

func debugChannel(s <-chan Work) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
