package vinna

import (
	"fmt"
	"time"
)

// Worker is the base interface for the workers in the worker pool
type Worker interface {
	Start()
	Stop()
}

// GenericWorker is a worker that will not post its results
type GenericWorker struct {
	ID      int
	jobs    <-chan Work
	sigStop chan bool
	working bool
}

// NewGenericWorker will create a new unstarted worker
func NewGenericWorker(n int, jobs <-chan Work) GenericWorker {
	return GenericWorker{
		ID:      n,
		jobs:    jobs,
		sigStop: make(chan bool),
		working: false,
	}
}

// Start will start the worker process for the GenericWorker
func (gw GenericWorker) Start() {
	if gw.working {
		return
	}
	go func() {
		fmt.Printf("Starting worker: \t%d\n", gw.ID)
		for {
			select {
			case w := <-gw.jobs:
				w.Do()
			case <-gw.sigStop:
				fmt.Printf("Stopping worker: \t%d\n", gw.ID)
				return
			default:
				time.Sleep(10 * time.Millisecond)
				fmt.Print(".")
			}
		}
	}()
}

// Stop will signal the worker to stop
func (gw *GenericWorker) Stop() {
	gw.working = false
	gw.sigStop <- true
}
