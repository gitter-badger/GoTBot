package twitch

import "github.com/thoj/go-ircevent"

type Handler interface {
	Handle(e *irc.Event)
}
