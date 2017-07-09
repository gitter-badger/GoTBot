package queue

import (
	"github.com/3stadt/GoTBot/structs"
	"github.com/3stadt/GoTBot/handlers"
	"github.com/thoj/go-ircevent"
	"fmt"
)

func commandWorker(job structs.Job, connection  *irc.Connection) {
	msg, err := handlers.CommandMap[job.Command](job.Channel, job.Sender, job.Params)
	if err != nil {
		fmt.Println(err)
	}
	connection.Privmsg(msg.Channel, msg.Message)
}

func HandleCommand(c chan structs.Job, connection *irc.Connection) {
	for {
		commandWorker(<- c, connection)
	}
}