package vinna

// Work is the base interface for the jobs to be done by our Workers.
// Any implementation that defines a self encapsulated Do() function
// can be passed into a jobs channel that will be done by the workers in the pool
type Work interface {
	Do()
}
