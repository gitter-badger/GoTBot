package handlers

import (
	"github.com/thoj/go-ircevent"
	"strings"
	"time"
	"math/rand"
	"github.com/3stadt/GoTBot/src/globals"
)

func Slap(channel string, sender string, params string, connection *irc.Connection) {
	victim := strings.TrimSpace(params)
	if len(params) < 1 || strings.ContainsAny(victim, " ") {
		return
	}

	if victim == "himself" || victim == "herself" || victim == "itself" || victim == globals.Conf["TWITCH_USER"] {
		connection.Privmsg(channel, "/me slaps "+sender+" playfully around with the mighty banhammer...")
		return
	}

	rand.Seed(time.Now().Unix())
	objects := []string{
		"a large trout",
		"no visible result",
		"the largest trout ever seen",
		"a barbie doll",
		"a blood stained sack",
		"a chainsaw",
	}
	n := rand.Int() % len(objects)
	connection.Privmsg(channel, sender+" slaps "+victim+" around a bit with "+objects[n]+"!")
}