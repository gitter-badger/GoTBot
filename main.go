package main

import (
	"github.com/thoj/go-ircevent"
	"crypto/tls"
	"fmt"
	"strings"
	"github.com/3stadt/GoTBot/src/handlers"
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/3stadt/GoTBot/src/bolt"
	"time"
	"github.com/joho/godotenv"
	"log"
)

var conf map[string]string
const serverSSL = "irc.chat.twitch.tv:443"

var commandMap = map[string]func(channel string, sender string, params string, connection *irc.Connection){
	"goSay":   handlers.Echo,
	"slap":    handlers.Slap,
}

func main() {
	var err error
	conf, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	channel := "#" + conf["TWITCH_CHANNEL"]
	botnick := conf["TWITCH_USER"]
	oauth := conf["OAUTH"]
	checkErr(err)
	oauthString := strings.TrimSpace(string(oauth))
	connection := irc.IRC(botnick, botnick)
	connection.VerboseCallbackHandler = true
	connection.Debug = true
	connection.UseTLS = true
	connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	connection.Password = oauthString

	connection.AddCallback("001", func(e *irc.Event) {
		connection.SendRaw("CAP REQ :twitch.tv/membership")
		connection.Join(channel)
	})
	connection.AddCallback("366", func(e *irc.Event) {})
	connection.AddCallback("PART", func(e *irc.Event) {
		if e.Nick == botnick {
			return
		}
		bolt.CreateOrUpdateUser(structs.User{
			Name: e.Nick,
			LastPart: time.Now(),
		})
	})
	connection.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick == botnick {
			return
		}
		bolt.CreateOrUpdateUser(structs.User{
			Name: e.Nick,
			LastJoin: time.Now(),
		})
	})

	connection.AddCallback("PRIVMSG", func(e *irc.Event) {
		if e.Nick == botnick {
			return
		}
		bolt.CreateOrUpdateUser(structs.User{
			Name: e.Nick,
			LastActive: time.Now(),
		})
		message := e.Message()
		if len(message) > 1 && strings.HasPrefix(message, "!") {
			i := strings.Index(message, " ")

			channel := e.Arguments[0]
			sender := e.Nick
			var command string
			var params string

			if i < 0 {
				command = message[1:]
				params = ""
			} else {
				command = message[1:i]
				params = message[i:]
			}

			if _, ok := commandMap[command]; ok {
				commandMap[command](channel, sender, params, connection)
			}

		}
	})

	err = connection.Connect(serverSSL)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	connection.Loop()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
