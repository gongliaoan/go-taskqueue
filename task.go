package taskqueue

import (
	"fmt"
	"time"
)

// TaskHandler define the task handler interface, any one that implements this interface can be appended to the queue
type TaskHandler interface {
	Success(id uint64)
	Timeout(id uint64)
}

// DebugTask struct for debug test object
type DebugTask struct {
}

// NewDebugTask create new DebugTask Object
func NewDebugTask() *DebugTask {
	return &DebugTask{}
}

// Success method will be called when the queue consumer take the task from the queue
func (t *DebugTask) Success(id uint64) {
	fmt.Println("# > running task:", id)
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Println("# < running task:", id)
}

// Timeout method will be called when the queue consumer take the task from the queue after timeout
func (t *DebugTask) Timeout(id uint64) {
	fmt.Println("# timeout:", id)
}
