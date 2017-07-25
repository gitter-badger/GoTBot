package handlers

import (
	"github.com/3stadt/GoTBot/src/structs"
)

func Echo(channel string, sender string, params string) (*structs.Message, error) {
	return &structs.Message{
		Channel: channel,
		Message: "Thanks's for sending goSay with '" + params + "' on " + channel + ", " + sender + "!",
	}, nil
}