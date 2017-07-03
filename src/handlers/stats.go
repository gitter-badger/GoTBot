package handlers

import (
	"github.com/thoj/go-ircevent"
	"strings"
	"github.com/3stadt/GoTBot/src/bolt"
)

func Stats(channel string, sender string, params string, connection *irc.Connection) {
	target := strings.TrimSpace(params)
	if len(params) < 1 || strings.ContainsAny(target, " ") {
		target = sender
	}
	targetUser := bolt.GetUser(target)
	if targetUser == nil {
		return
	}
	connection.Privmsg(channel, "User " + targetUser.Name + " was last active on " + targetUser.LastActive.Format("Mon, Jan 2 15:04:05"))
}
