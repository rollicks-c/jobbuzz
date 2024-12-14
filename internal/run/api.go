package run

import (
	"context"
	"fmt"
	"time"
)

type Job[T any] func() (T, error)

var TimeoutError = fmt.Errorf("timeout occured")

func WithTimeout[T any](timeout time.Duration, h Job[T]) (T, error) {

	// feedback channel
	doneChan := make(chan jobResult[T], 1)

	// content with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// execute
	go func() {
		val, err := h()
		res := jobResult[T]{
			res: val,
			err: err,
		}
		doneChan <- res
	}()

	// await completion
	var emptyRes T
	select {
	case <-ctx.Done():
		return emptyRes, TimeoutError
	case res := <-doneChan:
		return res.res, res.err
	}
}
