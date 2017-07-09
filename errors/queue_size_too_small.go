package fail

import "fmt"

type QueueSizeTooSmall struct {
	Min int
}

func (e *QueueSizeTooSmall) Error() string {
	return fmt.Sprintf("Queue size too small, min. %d required!", e.Min)
}
