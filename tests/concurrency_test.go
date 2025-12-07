package tests

import (
	"sync"
	"testing"

	"github.com/baranius/godiator"
	"github.com/stretchr/testify/assert"
)

type ConcRequest struct {
	ID int
}

type ConcResponse struct {
	ID int
}

type ConcHandler struct{}

func (h *ConcHandler) Handle(req ConcRequest, params ...any) (ConcResponse, error) {
	id := req.ID
	return ConcResponse{ID: id}, nil
}

func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	iterations := 100

	// Concurrent Registration and Sending
	for i := 0; i < iterations; i++ {
		wg.Add(2)

		// Routine 1: Register (Write)
		go func() {
			defer wg.Done()
			godiator.RegisterHandler[ConcRequest, ConcResponse](&ConcHandler{})
		}()

		// Routine 2: Send (Read)
		// Note: Send might fail if handler is not yet registered or temporarily replaced,
		// but we are testing that it DOES NOT PANIC due to map race.
		go func(val int) {
			defer wg.Done()
			_, _ = godiator.Send[ConcRequest, ConcResponse](ConcRequest{ID: val})
		}(i)
	}

	wg.Wait()
	assert.True(t, true, "Completed without panic")
}
