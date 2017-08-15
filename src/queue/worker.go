package queue

import (
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/3stadt/GoTBot/src/handlers"
	"github.com/thoj/go-ircevent"
	"fmt"
	"github.com/3stadt/GoTBot/src/db"
	"github.com/3stadt/GoTBot/src/res"
)

func commandWorker(job structs.Job, connection *irc.Connection, p *db.Pool, v *res.Vars) {
	var msg *structs.Message
	var err error
	if _, ok := handlers.CommandMap[job.Command]; ok {
		err = handlers.CommandMap[job.Command](job.Channel, job.Sender, job.Params, connection, p, v)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		if err := executeJsFile(job, connection, p, v); err != nil {
			fmt.Println(err)
		}
	}
	if msg != nil {
		connection.Privmsg(msg.Channel, msg.Message)
	}
}
func executeJsFile(job structs.Job, connection *irc.Connection, p *db.Pool, v *res.Vars) (error) {
	fileNames := handlers.PluginCommandMap[job.Command]
	for _, fileName := range fileNames {
		handlers.JsPlugin(fileName, job.Channel, job.Sender, job.Params, connection, p, v)
	}
	return nil
}

func HandleCommand(c chan structs.Job, connection *irc.Connection, p *db.Pool, v *res.Vars) {
	for {
		commandWorker(<-c, connection, p, v)
	}
}
