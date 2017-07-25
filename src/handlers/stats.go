package handlers

import (
	"strings"
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/3stadt/GoTBot/src/db"
	"strconv"
)

func Stats(channel string, sender string, params string) (*structs.Message, error) {
	target := strings.TrimSpace(params)
	if len(params) < 1 || strings.ContainsAny(target, " ") {
		target = sender
	}
	targetUser, err := db.GetUser(target)
	if err != nil {
		return nil, err
	}
	return &structs.Message{
		Channel: channel,
		Message: "User " + targetUser.Name + " with " + strconv.Itoa(targetUser.MessageCount) + " messages sent was last active on " + targetUser.LastActive.Format("Mon, Jan 2 15:04:05"),
	}, nil
}
