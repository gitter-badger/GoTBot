package handlers

import "github.com/thoj/go-ircevent"

func Echo(channel string, sender string, params string, connection *irc.Connection) {
	connection.Privmsg(channel, "Thanks's for sending goSay with '"+params+"' on "+channel+", "+sender+"!")
}