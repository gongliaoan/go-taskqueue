package taskqueue

import "fmt"

func consumer(messageCh <-chan *message, deleteCh chan<- Notification, queueID string) {

	for message := range messageCh {
		select {
		case <-message.cancelCh:
			message.taskHandler.Timeout(queueID, fmt.Sprint(message.id))
			close(message.timeoutCh)

		default:
			message.taskHandler.Success(queueID, fmt.Sprint(message.id))
			close(message.doneCh)
		}
	}

	close(deleteCh)
}
