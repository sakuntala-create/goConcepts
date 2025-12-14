package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

type Pool struct {
	m         sync.Mutex
	resources chan io.Closer
	factory   func() (io.Closer, error)
	closed    bool
}

var ErrPoolClosed = errors.New("pool has been closed")

func New(fn func() (io.Closer, error), size int) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size value too small")
	}
	return &Pool{
		resources: make(chan io.Closer, size),
		factory:   fn,
	}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources:
		if !ok {
			return nil, ErrPoolClosed
		}
		log.Printf("acquired existing resource")
		return r, nil
	default:
		// Provide a new resource since there are none available.
		log.Println("Acquire:", "New Resource")
		return p.factory()
	}
}

// Close will shutdown the pool and close all existing resources.
func (p *Pool) Close() {

	// Secure this operation with the Release operation.
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	p.closed = true

	close(p.resources)

	for r := range p.resources {
		r.Close()
	}

}

func (p *Pool) Release(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()
	if p.closed {
		r.Close()
		return
	}
	select {
	case p.resources <- r:
		log.Println("released resource")
	default:
		// Pool is full, close the resource.
		log.Println("release:", "closing resource")
		r.Close()
	}

}
