package taskqueue

// TaskHandler define the task handler interface, any one that implements this interface can be appended to the queue
type TaskHandler interface {
	Success(queueID, taskID string)
	Timeout(queueID, taskID string)
}
