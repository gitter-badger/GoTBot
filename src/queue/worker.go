package queue

import (
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/3stadt/GoTBot/src/handlers"
	"github.com/3stadt/GoTBot/src/globals"
)

func commandWorker(job structs.Job) {
	handlers.CommandMap[job.Command](job.Channel, job.Sender, job.Params, globals.Connection)
}

func HandleCommand(c chan structs.Job) {
	for {
		commandWorker(<- c)
	}
}