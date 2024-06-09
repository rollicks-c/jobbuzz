package pool

import (
	"context"
	"sync"
)

func Create[T any](size int, options ...Option[T]) *Pool[T] {
	// sanity check
	if size < 1 {
		panic("pool size must be at least 1")
	}

	// init pool
	p := &Pool[T]{
		data: &dataStore[T]{
			results: make([]T, 0),
			error:   nil,
		},
		size:            size,
		jobChan:         make(chan Job[T], 10000),
		resultChan:      make(chan jobResult[T]),
		resultsDoneChan: make(chan bool),
		wg:              &sync.WaitGroup{},
		setup:           &sync.Once{},
	}
	p.abortCtx, p.abortFunc = context.WithCancel(context.Background())
	p.shutdownCtx, p.shutdownFunc = context.WithCancel(context.Background())

	// apply op
	for _, option := range options {
		option(p)
	}
	return p
}

func (p Pool[T]) AddJob(job jobHandler[T], options ...JobOption[T]) {

	// ensure pool is started
	p.setup.Do(p.start)

	// increment job count
	p.data.jobCount++
	p.wg.Add(1)

	// add job to queue
	pJob := Job[T]{
		handler: job,
		JobTag:  "",
	}
	for _, option := range options {
		option(&pJob)
	}
	p.jobChan <- pJob
}

func (p Pool[T]) Complete() ([]T, error) {

	// ensure pool is /was running
	p.setup.Do(p.start)

	// wait for all jobs to complete or pool is aborted
	gracefulCompletionChan := p.gracefulCompletion()
	select {
	case <-p.abortCtx.Done():
	case <-gracefulCompletionChan:
	}

	// return results
	return p.data.results, p.data.error
}
