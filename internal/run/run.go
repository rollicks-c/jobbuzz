package run

type jobResult[T any] struct {
	res T
	err error
}
