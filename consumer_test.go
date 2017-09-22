package taskqueue

import (
	"testing"
	"time"

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

	go consumer(messageChannel, deleteChannel, queueID)

	_, _, msg := newMessage(time.Duration(0), &taskTestConsumerSendMessage{})
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
	queueID string
	t       *testing.T
}

func (task *taskTestConsumerSendTask) Success(queueID, taskID string) {

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
		queueID: "consumer",
	}

	go consumer(messageChannel, deleteChannel, task.queueID)

	doneOutCh, timeoutOutCh, msg := newMessage(time.Duration(0), task)
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
