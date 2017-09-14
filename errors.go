package taskqueue

import "errors"

var (
	// ErrTaskQueueFull error object fired when queue are full, cap == len
	ErrTaskQueueFull = errors.New("message channel is full")
)
