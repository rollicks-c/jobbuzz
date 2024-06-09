package pool

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestComplete(t *testing.T) {

	pl := Create[any](5)

	for i := 0; i < 10; i++ {

		pl.AddJob(func() (any, error) {

			time.Sleep(1 * time.Second)
			return nil, nil

		})

	}

	_, err := pl.Complete()
	assert.NoError(t, err)

}

func TestAbort(t *testing.T) {

	pl := Create[any](5)

	for i := 0; i < 10; i++ {

		pl.AddJob(func() (any, error) {

			if i == 7 {
				return nil, fmt.Errorf("error")
			}

			time.Sleep(1 * time.Second)

			return nil, nil

		})

	}

	_, err := pl.Complete()
	assert.Error(t, err)

}
