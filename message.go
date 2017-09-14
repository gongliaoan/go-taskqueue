package taskqueue

import "time"

type message struct {
	doneCh, timeoutCh chan<- Notification // output channel
	cancelCh          <-chan Notification // input channel

	id          uint64
	timer       *time.Timer
	taskHandler TaskHandler
}

func newMessage(id uint64, timeout time.Duration, taskHandler TaskHandler) (doneCh, timeoutCh <-chan Notification, msg *message) {

	cancelChannel := make(chan Notification)
	doneChannel := make(chan Notification)
	timeoutChannel := make(chan Notification)

	msg = func(cancelCh <-chan Notification, doneCh, timeoutCh chan<- Notification) *message {
		return &message{
			cancelCh:    cancelCh,
			doneCh:      doneCh,
			timeoutCh:   timeoutCh,
			taskHandler: taskHandler,
			id:          id,
			timer:       nil,
		}
	}(cancelChannel, doneChannel, timeoutChannel)

	msg.timer = func(cancelCh chan<- Notification, timeout time.Duration) *time.Timer {
		return time.AfterFunc(timeout, func() {
			close(cancelCh)
		})
	}(cancelChannel, timeout)

	return doneChannel, timeoutChannel, msg
}
