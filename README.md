# Vinna - Self Adjusting Worker Pool
---
Vinna is a Go implementation of a work queue with a manager for the worker pool. The manager controls the number of workers in the pool by monitoring the volume of the Work channel in comparison to the maximum capacity of the channel.

In order to use the worker pool it is necessary to implement the Work interface and pass the channel of Work to the WorkerPool

The current version will both add and remove workers based on the volume of work in the channel. A worker wil be added if the volume of the channel buffer is above 20% of the total capacity. A worker will be removed if the volume drops below 10% of the total capacity. In future versions most of this will be configurable.

#### Worker Interface
```go
// Work is the base interface for the jobs to be done by our Workers.
type Work interface {
    Do()
}
```



#### Example
```go
type TestWork struct {
	Msg string	
}
func (w TestWork) Do() {
	fmt.Print("*")
	return
}

func main() {
	jobs := make(chan Work, 10)

	// The number of workers we want to have initially
	n := 1
	pool := NewWorkerPool(n, jobs)

	// This will start our worker pool
	pool.StartAllWorkers()

	// Interrupt handler to clean up when running in command line
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	
	// Simulate putting work to channel
	go func() {
		for i := 1; i <= 20; i++ {
        	jobs <- TestWork{Msg: "test"}
    	}
	}()
	
	// Block to allow goroutines to run indefinitely 
	<-c
    pool.Stop()
    os.Exit(0)
}
```
