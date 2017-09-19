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

// begin FullQueue
type testFullQueue struct {
	wait   chan Notification
	notify chan Notification
}

func (task *testFullQueue) Success(queueID, taskID string) {
	close(task.notify)
	<-task.wait
}

func (task *testFullQueue) Timeout(queueID, taskID string) {}

func Test_QueueFullAsync(t *testing.T) {

	queue := New("FullQueueAsync", 1, time.Duration(0))
	defer queue.Close()

	var err [3]error
	var task [3]*testFullQueue
	var doneCh, timeoutCh [3]<-chan Notification

	for i := 0; i < 3; i++ {
		task[i] = &testFullQueue{
			wait:   make(chan Notification),
			notify: make(chan Notification),
		}
	}

	doneCh[0], timeoutCh[0], err[0] = queue.EnqueueAsync(task[0])
	<-task[0].notify
	doneCh[1], timeoutCh[1], err[1] = queue.EnqueueAsync(task[1])
	doneCh[2], timeoutCh[2], err[2] = queue.EnqueueAsync(task[2])

	require.NoError(t, err[0])
	require.NoError(t, err[1])
	require.Error(t, err[2])

	close(task[0].wait)
	<-task[1].notify
	close(task[1].wait)

	for i := 0; i < 2; i++ {
		var success bool
		select {
		case <-doneCh[i]:
			success = true
		case <-timeoutCh[i]:
			success = false
		}
		require.True(t, success)
	}
}

func Test_QueueFullSync(t *testing.T) {

	queue := New("FullQueueSync", 1, time.Duration(0))
	defer queue.Close()

	var err [3]error
	var task [3]*testFullQueue
	var doneCh, timeoutCh [2]<-chan Notification

	for i := 0; i < 3; i++ {
		task[i] = &testFullQueue{
			wait:   make(chan Notification),
			notify: make(chan Notification),
		}
	}

	doneCh[0], timeoutCh[0], err[0] = queue.EnqueueAsync(task[0])
	<-task[0].notify
	doneCh[1], timeoutCh[0], err[1] = queue.EnqueueAsync(task[1])
	err[2] = queue.Enqueue(task[2])

	require.NoError(t, err[0])
	require.NoError(t, err[1])
	require.Error(t, err[2])

	close(task[0].wait)
	<-task[1].notify
	close(task[1].wait)

	for i := 0; i < 2; i++ {
		var success bool
		select {
		case <-doneCh[i]:
			success = true
		case <-timeoutCh[i]:
			success = false
		}
		require.True(t, success)
	}
}

// end FullQueue
