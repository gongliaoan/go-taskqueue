package taskqueue

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testTask struct {
}

func (t *testTask) Success(queueID, taskID string) {
	fmt.Println("# > running task:", taskID, "from queue:", queueID)
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Println("# < running task:", taskID, "from queue:", queueID)
}

func (t *testTask) Timeout(queueID, taskID string) {
	fmt.Println("# timeout:", taskID, "from queue:", queueID)
}

func TestNew(t *testing.T) {
	queue := New("test-queue", 10, time.Duration(1))
	defer queue.Close()

	require.NotNil(t, queue)
	require.False(t, isFull(queue))
}

func TestQueue_NewTask(t *testing.T) {

	var (
		numberOfTasks    = 10
		secondsToTimeout = time.Duration(5) * time.Second
	)

	queue := New("test-queue", numberOfTasks, secondsToTimeout)
	defer queue.Close()

	wg := sync.WaitGroup{}
	wg.Add(numberOfTasks)

	for i := 0; i < numberOfTasks; i++ {
		go func() {
			defer wg.Done()
			err := queue.Enqueue(&testTask{})
			require.NoError(t, err)
		}()
	}

	wg.Wait()
}
