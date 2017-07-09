package handlers

import (
	"strings"
	"github.com/3stadt/GoTBot/bolt"
	"github.com/3stadt/GoTBot/structs"
	"github.com/3stadt/GoTBot/errors"
)

func Stats(channel string, sender string, params string) (*structs.Message, error) {
	target := strings.TrimSpace(params)
	if len(params) < 1 || strings.ContainsAny(target, " ") {
		target = sender
	}
	targetUser := bolt.GetUser(target)
	if targetUser == nil {
		return nil, &fail.NoTargetUser{Name: target}
	}
	return &structs.Message{
		Channel: channel,
		Message: "User " + targetUser.Name + " was last active on " + targetUser.LastActive.Format("Mon, Jan 2 15:04:05"),
	}, nil
}
