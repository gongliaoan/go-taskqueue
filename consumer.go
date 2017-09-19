package taskqueue

import "strconv"

func consumer(messageCh <-chan *message, deleteCh chan<- Notification, queueID string) {

	for message := range messageCh {
		select {
		case <-message.cancelCh:
			message.taskHandler.Timeout(queueID, strconv.FormatUint(message.id, 10))
			close(message.timeoutCh)

		default:
			message.taskHandler.Success(queueID, strconv.FormatUint(message.id, 10))
			close(message.doneCh)
		}
	}

	close(deleteCh)
}
