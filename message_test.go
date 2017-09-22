package taskqueue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type taskTestMessage struct{}

func (t *taskTestMessage) Success(queueID, taskID string) {}
func (t *taskTestMessage) Timeout(queueID, taskID string) {}

func Test_messageNoTimeout(t *testing.T) {

	doneOutCh, timeOutCh, msg := newMessage(time.Duration(0), &taskTestMessage{})
	require.NotNil(t, msg)
	require.Nil(t, msg.timer)

	done := false
	lockDone := make(chan Notification)
	go func() {
		<-doneOutCh
		done = true
		close(lockDone)
	}()
	close(msg.doneCh)
	<-lockDone
	require.True(t, done)

	timeout := false
	lockTimeout := make(chan Notification)
	go func() {
		<-timeOutCh
		timeout = true
		close(lockTimeout)
	}()
	close(msg.timeoutCh)
	<-lockTimeout
	require.True(t, timeout)

	cancel := false
	time.Sleep(time.Duration(5) * time.Millisecond)
	select {
	case <-msg.cancelCh:
		cancel = true
	default:
	}
	require.False(t, cancel)
}

func Test_messageWithTimeout(t *testing.T) {

	doneOutCh, timeOutCh, msg := newMessage(time.Duration(10)*time.Millisecond, &taskTestMessage{})
	require.NotNil(t, msg)
	require.NotNil(t, msg.timer)

	done := false
	lockDone := make(chan Notification)
	go func() {
		<-doneOutCh
		done = true
		close(lockDone)
	}()
	close(msg.doneCh)
	<-lockDone
	require.True(t, done)

	timeout := false
	lockTimeout := make(chan Notification)
	go func() {
		<-timeOutCh
		timeout = true
		close(lockTimeout)
	}()
	close(msg.timeoutCh)
	<-lockTimeout
	require.True(t, timeout)

	cancel := false
	select {
	case <-msg.cancelCh:
		cancel = true
	default:
	}
	require.False(t, cancel)

	time.Sleep(time.Duration(20) * time.Millisecond)
	select {
	case <-msg.cancelCh:
		cancel = true
	default:
	}
	require.True(t, cancel)
}

// begin message timeout
type taskTestMessageTimeout struct {
	timeout bool
}

func (task *taskTestMessageTimeout) Success(queueID, taskID string) {
	task.timeout = false
}

func (task *taskTestMessageTimeout) Timeout(queueID, taskID string) {
	task.timeout = true
}

func Test_messageTimeout(t *testing.T) {

	var (
		messageChannelSend chan<- *message
		deleteChannelRecv  <-chan Notification
	)

	messageChannel := make(chan *message, 1)
	deleteChannel := make(chan Notification)

	messageChannelSend = messageChannel
	deleteChannelRecv = deleteChannel

	// force timeout by hand
	timeoutHandle := func(cancelCh chan<- Notification, timeout time.Duration) *time.Timer {
		close(cancelCh)
		return nil
	}

	task := &taskTestMessageTimeout{
		timeout: false,
	}

	doneCh, timeoutCh, msg := newMessageWithTimeoutHandleFunc(time.Duration(1), task, timeoutHandle)
	require.NotNil(t, msg)
	require.Nil(t, msg.timer)

	go consumer(messageChannel, deleteChannel, "messageTimeout")

	messageChannelSend <- msg
	close(messageChannel)
	<-deleteChannelRecv

	var timeout bool
	select {
	case <-doneCh:
		timeout = false
	case <-timeoutCh:
		timeout = true
	}

	require.True(t, timeout)
	require.True(t, task.timeout)
}

// end message timeout
