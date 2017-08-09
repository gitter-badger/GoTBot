package twitch

import (
	"github.com/thoj/go-ircevent"
	"github.com/3stadt/GoTBot/src/context"
)

func (c *Client) Connected(e *irc.Event) {
	channel := "#" + context.Conf["TWITCH_CHANNEL"]
	c.Connection.SendRaw("CAP REQ :twitch.tv/membership")
	c.Connection.SendRaw("CAP REQ :twitch.tv/commands")
	c.Connection.Join(channel)
	c.Connection.Privmsg(channel, "/mods")
}
