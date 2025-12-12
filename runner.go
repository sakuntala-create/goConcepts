// Package runner manages the running and lifetime of a process.
package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

type Runner struct {
	interrupt chan os.Signal
	complete  chan error
	timeout   <-chan time.Time
	// tasks holds a set of functions that are executed
	// synchronously in index order.
	tasks []func(int)
}

// ErrTimeout is returned when a value is received on the timeout.
var ErrTimeout = errors.New("received timeout")

// ErrInterrupt is returned when an event from the OS is received.
var ErrInterrupt = errors.New("received interrupt")

// New creates a new Runner.

func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

func(r *Runner) Add(tasks ...func(int)){
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) Start() error {

	signal.Notify(r.interrupt, os.Interrupt)
	go func(){
		r.complete <- r.run()
	}()
select {
	case err := <-r.complete:
		return err
	case <-r.timeout:
		return ErrTimeout
	}
}

func( r *Runner) run() error {
	for id, task := range r.tasks {
		if r.gotInterrupt(){
			return ErrInterrupt
		}
		task(id)
	}
	return nil
}

func (r *Runner) gotInterrupt() bool {
	select {
	case <-r.interrupt:	
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}
