package taskqueue

import "fmt"

func consumer(messageCh <-chan *message, deleteCh chan<- Notification, queueID uint64) {

	for message := range messageCh {
		select {
		case <-message.cancelCh:
			message.taskHandler.Timeout(message.id)
			close(message.timeoutCh)

		default:
			fmt.Println("- default, queue:", queueID, "message:", message.id)
			message.taskHandler.Success(message.id)
			close(message.doneCh)
		}
		fmt.Println("- loop, queue:", queueID, "message:", message.id, "len:", len(messageCh))
	}

	fmt.Println("- Exiting consumer for queue:", queueID)
	close(deleteCh)
}
