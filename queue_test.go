package taskqueue

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testDebugTask struct {
}

func (t *testDebugTask) Success(queueID, taskID string) {
	fmt.Println("# > running task:", taskID, "from:", queueID)
	time.Sleep(time.Duration(2) * time.Millisecond)
	fmt.Println("# < running task:", taskID, "from:", queueID)
}

func (t *testDebugTask) Timeout(queueID, taskID string) {
	fmt.Println("# timeout:", taskID, "from:", queueID)
}

func Test_QueueNew(t *testing.T) {
	queue := New("test-queue", 10, time.Duration(1))
	defer queue.Close()

	require.NotNil(t, queue)
	require.False(t, isFull(queue))
}

func Test_QueueNewTask(t *testing.T) {

	var (
		numberOfTasks    = 10
		secondsToTimeout = time.Duration(5) * time.Millisecond
	)

	queue := New("NewTask", numberOfTasks, secondsToTimeout)
	defer queue.Close()

	wg := sync.WaitGroup{}
	wg.Add(numberOfTasks)

	for i := 0; i < numberOfTasks; i++ {
		go func() {
			defer wg.Done()
			err := queue.Enqueue(&testDebugTask{})
			require.NoError(t, err)
		}()
	}

	wg.Wait()
}

// begin NoTimeout
type testTaskNoTimeout struct {
	successCount int64
	timeoutCount int64
}

func (t *testTaskNoTimeout) Success(queueID, taskID string) {
	atomic.AddInt64(&t.successCount, 1)
}

func (t *testTaskNoTimeout) Timeout(queueID, taskID string) {
	atomic.AddInt64(&t.timeoutCount, 1)
}

func Test_QueueNoTimeout(t *testing.T) {

	numberOfTasks := 10

	queue := New("NoTimeout", numberOfTasks, time.Duration(0))
	defer queue.Close()

	wg := sync.WaitGroup{}
	wg.Add(int(numberOfTasks))

	task := &testTaskNoTimeout{}

	for i := 0; i < numberOfTasks; i++ {
		go func() {
			defer wg.Done()
			err := queue.Enqueue(task)
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	require.Equal(t, int64(numberOfTasks), task.successCount)
	require.Equal(t, int64(0), task.timeoutCount)
}

// end NoTimeout

// begin NoSleep
type testTaskNoSleep struct {
	successCount int64
	timeoutCount int64
}

func (t *testTaskNoSleep) Success(queueID, taskID string) {
	atomic.AddInt64(&t.successCount, 1)
}

func (t *testTaskNoSleep) Timeout(queueID, taskID string) {
	atomic.AddInt64(&t.timeoutCount, 1)
}

func Test_QueueNoSleep(t *testing.T) {

	numberOfTasks := 10

	queue := New("NoSleep", numberOfTasks, time.Duration(1)*time.Second)
	defer queue.Close()

	wg := sync.WaitGroup{}
	wg.Add(int(numberOfTasks))

	task := &testTaskNoSleep{}

	for i := 0; i < numberOfTasks; i++ {
		go func() {
			defer wg.Done()
			err := queue.Enqueue(task)
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	require.Equal(t, int64(numberOfTasks), task.successCount+task.timeoutCount)
}

// end NoSleep
