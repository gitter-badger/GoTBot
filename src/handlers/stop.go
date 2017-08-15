package handlers

import (
	"errors"
	"strings"
)

func (d *deps) Stop() error {
	if d.sender == strings.TrimPrefix(d.channel, "#") || d.v.IsTwitchMod(d.sender) {
		d.connection.Privmsg(d.channel, "Shutting down bot...")
		d.connection.Quit()
		return nil
	}
	return errors.New("Insufficient permissions")
}
