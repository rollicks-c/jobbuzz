package pool

import (
	"context"
	"sync"
)

type Pool[T any] struct {

	// config
	Subscriber ProgressSubscriber

	// data
	data *dataStore[T]

	// setup
	size  int
	setup *sync.Once

	// concurrency work
	jobChan         chan Job[T]
	resultChan      chan jobResult[T]
	resultsDoneChan chan bool

	// concurrency control
	abortCtx     context.Context
	abortFunc    context.CancelFunc
	shutdownCtx  context.Context
	shutdownFunc context.CancelFunc
	wg           *sync.WaitGroup
}

type dataStore[T any] struct {
	jobCount int
	results  []T
	error    error
}

type jobResult[T any] struct {
	jobTag string
	result T
	err    error
}

func (p Pool[T]) start() {
	go p.processResults()
	for i := 0; i < p.size; i++ {
		go p.worker()
	}
}

func (p Pool[T]) processResults() {

	for {

		// listen for results
		jr, ok := <-p.resultChan
		if !ok {
			break
		}

		// abort if error
		if jr.err != nil {
			p.data.error = jr.err
			p.abort()
			return
		}

		// collect result
		p.data.results = append(p.data.results, jr.result)

		// report progress
		if p.Subscriber == nil {
			continue
		}
		p.Subscriber(jr.jobTag, len(p.data.results), p.data.jobCount)
	}

	// notify graceful shutdown
	p.resultsDoneChan <- true
}

func (p Pool[T]) worker() {
	for {

		// listen for abort or jobs
		select {
		case <-p.shutdownCtx.Done():
			return
		case job, ok := <-p.jobChan:

			// channel closed -> abort ops
			if !ok {
				p.wg.Done()
				return
			}

			// execute job
			data, err := job.handler()

			// send result
			result := jobResult[T]{
				jobTag: job.JobTag,
				result: data,
				err:    err,
			}
			p.resultChan <- result
			p.wg.Done()
		}
	}
}

func (p Pool[T]) gracefulCompletion() chan bool {

	// install notification channel
	completionChan := make(chan bool)

	// wait for all jobs to complete
	go func() {

		// wait job completion
		p.wg.Wait()

		// stop workers
		p.shutdownFunc()

		// wait for result processing to complete
		close(p.resultChan)
		<-p.resultsDoneChan
		completionChan <- true
	}()

	return completionChan

}

func (p Pool[T]) abort() {
	p.abortFunc()
	close(p.jobChan)
}
