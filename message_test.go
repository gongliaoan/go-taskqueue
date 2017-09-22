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

	messageID := uint64(1)
	doneOutCh, timeOutCh, msg := newMessage(messageID, time.Duration(0), &taskTestMessage{})
	require.NotNil(t, msg)
	require.Equal(t, msg.id, messageID)
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

	messageID := uint64(1)
	doneOutCh, timeOutCh, msg := newMessage(messageID, time.Duration(10)*time.Millisecond, &taskTestMessage{})
	require.NotNil(t, msg)
	require.Equal(t, msg.id, messageID)
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
