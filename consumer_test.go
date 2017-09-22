package taskqueue

import (
	"testing"
	"time"

	"strconv"

	"github.com/stretchr/testify/require"
)

// begin SendMessage
type taskTestConsumerSendMessage struct{}

func (t *taskTestConsumerSendMessage) Success(queueID, taskID string) {}
func (t *taskTestConsumerSendMessage) Timeout(queueID, taskID string) {}

func Test_ConsumerSendMessage(t *testing.T) {

	messageChannel := make(chan *message, 10)
	deleteChannel := make(chan Notification)
	deleted := false

	queueID := "consumer"
	messageID := uint64(0)

	go consumer(messageChannel, deleteChannel, queueID)

	_, _, msg := newMessage(messageID, time.Duration(0), &taskTestConsumerSendMessage{})
	messageChannel <- msg
	close(messageChannel)

	lockDelete := make(chan Notification)
	go func() {
		<-deleteChannel
		deleted = true
		close(lockDelete)
	}()

	<-lockDelete
	require.True(t, deleted)
}

// end SendMessage

// begin SendTask
type taskTestConsumerSendTask struct {
	queueID   string
	messageID uint64
	t         *testing.T
}

func (task *taskTestConsumerSendTask) Success(queueID, taskID string) {

	messageID, err := strconv.ParseUint(taskID, 10, 64)
	require.NoError(task.t, err)
	require.Equal(task.t, task.messageID, messageID)
	require.Equal(task.t, task.queueID, queueID)

}
func (task *taskTestConsumerSendTask) Timeout(queueID, taskID string) {
	require.Fail(task.t, "Unnexpected timeout :(")
}

func Test_ConsumerSendTask(t *testing.T) {

	messageChannel := make(chan *message, 1)
	deleteChannel := make(chan Notification)
	deleted := false

	task := &taskTestConsumerSendTask{
		queueID:   "consumer",
		messageID: 0,
	}

	go consumer(messageChannel, deleteChannel, task.queueID)

	doneOutCh, timeoutOutCh, msg := newMessage(task.messageID, time.Duration(0), task)
	messageChannel <- msg
	select {
	case <-doneOutCh:
	case <-timeoutOutCh:
		require.Fail(t, "Unnexpected timeout :(")
	}

	close(messageChannel)

	lockDelete := make(chan Notification)
	go func() {
		<-deleteChannel
		deleted = true
		close(lockDelete)
	}()

	<-lockDelete
	require.True(t, deleted)
}

// end SendTask
