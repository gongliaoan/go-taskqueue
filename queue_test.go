package taskqueue

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testTask struct {
}

func (t *testTask) Success(queueID, taskID string) {
	fmt.Println("# > running task:", taskID, "from:", queueID)
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Println("# < running task:", taskID, "from:", queueID)
}

func (t *testTask) Timeout(queueID, taskID string) {
	fmt.Println("# timeout:", taskID, "from:", queueID)
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

	queue := New("NewTask", numberOfTasks, secondsToTimeout)
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

func TestQueue_NoTimeout(t *testing.T) {

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

func TestQueue_NoSleep(t *testing.T) {

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
