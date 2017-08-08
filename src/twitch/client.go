package twitch

import "github.com/thoj/go-ircevent"

type Client struct {
	Connection *irc.Connection
	Nick string
	Moderators []string
	CommandQueueName string
}

func Init(oauth string, nick string, commandQueueName string) Client {
	client := Client{
		Connection: irc.IRC(nick, nick),
		Nick: nick,
	}
	client.Connection.Password = oauth
	client.CommandQueueName = commandQueueName
	client.Moderators = []string{}
	return client
}

