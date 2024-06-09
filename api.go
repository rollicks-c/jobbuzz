package jobbuzz

import (
	"github.com/rollicks-c/jobbuzz/internal/pool"
)

func WithTag[T any](tag string) pool.JobOption[T] {
	return func(job *pool.Job[T]) {
		job.JobTag = tag
	}
}

func WithProgress[T any](subscriber pool.ProgressSubscriber) pool.Option[T] {
	return func(p *pool.Pool[T]) {
		p.Subscriber = subscriber
	}
}

func New[T any](size int, options ...pool.Option[T]) *pool.Pool[T] {
	return pool.Create(size, options...)
}
