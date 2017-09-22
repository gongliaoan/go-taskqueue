package taskqueue

import (
	"time"
)

// Queue struct for queue object/class
type Queue struct {
	messageCh      chan<- *message
	deleteCh       <-chan Notification
	messageTimeout time.Duration
	queueID        string
}

// New queue constructor
func New(id string, cap int, timeout time.Duration) *Queue {

	messageChannel := make(chan *message, cap)
	deleteChannel := make(chan Notification)

	go consumer(messageChannel, deleteChannel, id)

	return &Queue{
		messageCh:      messageChannel,
		deleteCh:       deleteChannel,
		messageTimeout: timeout,
		queueID:        id,
	}
}

// CloseAsync send notification to queue deleted and returns a read only channel to user receive a Notification when
// deletion be completed
func (q *Queue) CloseAsync() <-chan Notification {
	close(q.messageCh)
	return q.deleteCh
}

// Close wait for all tasks be completed, after that, kill the consumer
func (q *Queue) Close() {
	<-q.CloseAsync()
}

// EnqueueAsync send a TaskHandler to the queue and return notification channels
func (q *Queue) EnqueueAsync(taskHandler TaskHandler) (doneCh, timeoutCh <-chan Notification, err error) {

	doneCh, timeoutCh, message := newMessage(q.messageTimeout, taskHandler)

	select {
	case q.messageCh <- message:
		return doneCh, timeoutCh, nil
	default:
		return nil, nil, ErrTaskQueueFull
	}
}

// Enqueue send a TaskHandler to the queue and wait for the task execution or timeout
func (q *Queue) Enqueue(taskHandler TaskHandler) (err error) {

	var doneCh, timeoutCh <-chan Notification

	if doneCh, timeoutCh, err = q.EnqueueAsync(taskHandler); err != nil {
		return err
	}

	select {
	case <-doneCh:
		return nil
	case <-timeoutCh:
		return nil
	}
}
