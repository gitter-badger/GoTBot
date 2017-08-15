package handlers

import (
	"github.com/thoj/go-ircevent"
	"errors"
	"strings"
	"github.com/3stadt/GoTBot/src/res"
	"github.com/3stadt/GoTBot/src/db"
)

func Stop(channel string, sender string, params string, connection *irc.Connection, p *db.Pool, v *res.Vars) error {
	if sender == strings.TrimPrefix(channel, "#") || v.IsTwitchMod(sender) {
		connection.Privmsg(channel, "Shutting down bot...")
		connection.Quit()
		return nil
	}
	return errors.New("Insufficient permissions")
}
