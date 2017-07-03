package queue

import (
	"errors"
	"github.com/3stadt/GoTBot/src/structs"
)

var JobChannels = make(map[string](chan structs.Job))

func NewQueue(name string, maxJobs int) error {
	var err error
	if maxJobs < 1 {
		maxJobs = 1
		err = errors.New("maxJobs must be at least 1. maxJobs was set to 1.")
	}
	JobChannels[name] = make(chan structs.Job, maxJobs)
	return err
}

func AddJob(queue string, job structs.Job) {
	JobChannels[queue] <- job
}