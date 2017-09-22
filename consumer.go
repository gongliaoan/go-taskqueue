package taskqueue

import "strconv"

func consumer(messageCh <-chan *message, deleteCh chan<- Notification, queueID string) {

	messageID := uint64(0)

	for message := range messageCh {
		select {
		case <-message.cancelCh:
			message.taskHandler.Timeout(queueID, strconv.FormatUint(messageID, 10))
			close(message.timeoutCh)

		default:
			message.taskHandler.Success(queueID, strconv.FormatUint(messageID, 10))
			close(message.doneCh)
		}

		messageID++
	}

	close(deleteCh)
}
