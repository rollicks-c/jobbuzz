package pool

type ProgressSubscriber func(jobTag string, current, total int)

type JobOption[T any] func(*Job[T])

type Option[T any] func(*Pool[T])

type Job[T any] struct {
	handler jobHandler[T]
	JobTag  string
}

type jobHandler[T any] func() (T, error)
