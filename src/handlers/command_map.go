package handlers

import (
	"github.com/thoj/go-ircevent"
	"github.com/3stadt/GoTBot/src/db"
	"github.com/3stadt/GoTBot/src/res"
)

var CommandMap = map[string]func(channel string, sender string, params string, connection *irc.Connection, p *db.Pool, v *res.Vars) error{
	"goSay":    Echo,
	"slap":     Slap,
	"stats":    Stats,
	"shutdown": Stop,
	"stop":     Stop,
}
