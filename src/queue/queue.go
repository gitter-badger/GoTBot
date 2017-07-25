package queue

import (
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/3stadt/GoTBot/src/errors"
)

var JobQueue = make(map[string](chan structs.Job))

func NewQueue(name string, maxJobs int) error {
	var err error
	if maxJobs < 1 {
		maxJobs = 1
		err = &fail.QueueSizeTooSmall{Min: 1}
	}
	JobQueue[name] = make(chan structs.Job, maxJobs)
	return err
}

func AddJob(queue string, job structs.Job) {
	JobQueue[queue] <- job
}