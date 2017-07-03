package handlers

import (
	"github.com/thoj/go-ircevent"
	"strings"
	"github.com/3stadt/GoTBot/src/bolt"
	"fmt"
	"encoding/json"
)

func Stats(channel string, sender string, params string, connection *irc.Connection) {
	target := strings.TrimSpace(params)
	if len(params) < 1 || strings.ContainsAny(target, " ") {
		target = sender
	}
	targetUser := bolt.GetUser(target)
	if targetUser == nil {
		fmt.Println(targetUser, targetUser)
		return
	}
	targetData, err := json.Marshal(targetUser)
	if err != nil {
		fmt.Println(targetData, err)
		return
	}
	connection.Privmsg(channel, string(targetData))
}
