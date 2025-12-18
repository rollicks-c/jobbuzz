package run

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSuccess(t *testing.T) {
	job := func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "success", nil
	}
	result, err := WithTimeout(200*time.Millisecond, job)
	assert.NoError(t, err)
	assert.Equal(t, "success", result)
}

func TestTimeout(t *testing.T) {
	job := func() (string, error) {
		time.Sleep(300 * time.Millisecond)
		return "success", nil
	}
	result, err := WithTimeout(200*time.Millisecond, job)
	assert.Error(t, err)
	assert.True(t, errors.Is(TimeoutError, err))
	assert.Empty(t, result)
}

func TestJobFail(t *testing.T) {
	job := func() (string, error) {
		return "", errors.New("job error")
	}
	result, err := WithTimeout(200*time.Millisecond, job)
	assert.Error(t, err)
	assert.Equal(t, "job error", err.Error())
	assert.Empty(t, result)
}
