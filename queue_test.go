package taskqueue

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	queue := NewQueue("test-queue", 10, time.Duration(1))
	defer queue.Close()

	require.NotNil(t, queue)
	require.False(t, isFull(queue))
}

func TestQueue_NewTask(t *testing.T) {

	var (
		numberOfTasks    = 10
		secondsToTimeout = time.Duration(5) * time.Second
	)

	queue := NewQueue("test-queue", numberOfTasks, secondsToTimeout)
	defer queue.Close()

	wg := sync.WaitGroup{}
	wg.Add(numberOfTasks)

	for i := 0; i < numberOfTasks; i++ {
		go func() {
			defer wg.Done()
			err := queue.Enqueue(NewDebugTask())
			require.NoError(t, err)
		}()
	}

	wg.Wait()
}
