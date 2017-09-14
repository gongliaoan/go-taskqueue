package tq

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	queue := NewQueue(0, 10, time.Duration(1))
	defer queue.Delete()

	require.NotNil(t, queue)
	require.False(t, isFull(queue))
}

func TestQueue_NewTask(t *testing.T) {

	var (
		numberOfTasks    = 10
		secondsToTimeout = time.Duration(5) * time.Second
	)

	queue := NewQueue(0, numberOfTasks, secondsToTimeout)
	defer queue.Delete()

	wg := sync.WaitGroup{}
	wg.Add(numberOfTasks)

	for i := 0; i < numberOfTasks; i++ {
		go func() {
			defer wg.Done()
			err := queue.Append(NewDebugTask())
			require.NoError(t, err)
		}()
	}

	wg.Wait()
}
