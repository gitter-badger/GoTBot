package handlers

import (
	"github.com/thoj/go-ircevent"
	"github.com/3stadt/GoTBot/src/db"
	"github.com/3stadt/GoTBot/src/res"
)

func Echo(channel string, sender string, params string, connection *irc.Connection, p *db.Pool, v *res.Vars) error {
	connection.Privmsg(channel, "Thanks's for sending goSay with '"+params+"' on "+channel+", "+sender+"!")
	return nil
}
