package twitch

import (
	"github.com/thoj/go-ircevent"
	"strings"
)

func (c *Client) Notice(e *irc.Event) {
	message := e.Message()
	moderatorPrefix := "The moderators of this room are: "
	if strings.HasPrefix(message, moderatorPrefix) {
		c.Moderators = strings.Split(strings.TrimPrefix(message, moderatorPrefix), ", ")
	}
}
