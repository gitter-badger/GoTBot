package main

import (
	"github.com/thoj/go-ircevent"
	"crypto/tls"
	"fmt"
	"strings"
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/3stadt/GoTBot/src/bolt"
	"time"
	"github.com/joho/godotenv"
	"log"
	"github.com/3stadt/GoTBot/src/queue"
	"github.com/3stadt/GoTBot/src/handlers"
	"github.com/3stadt/GoTBot/src/globals"
	"strconv"
)

const serverSSL = "irc.chat.twitch.tv:443"

func main() {
	var err error
	globals.Conf, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	queue.NewQueue(globals.CommandQueueName, 30)
	channel := "#" + globals.Conf["TWITCH_CHANNEL"]
	botnick := globals.Conf["TWITCH_USER"]
	oauth := globals.Conf["OAUTH"]
	debug, debugErr := strconv.ParseBool(globals.Conf["DEBUG"])
	if debugErr != nil {
		debug = false
	}
	checkErr(err)
	oauthString := strings.TrimSpace(string(oauth))
	globals.Connection = irc.IRC(botnick, botnick)
	globals.Connection.VerboseCallbackHandler = debug
	globals.Connection.Debug = debug
	globals.Connection.UseTLS = true
	globals.Connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	globals.Connection.Password = oauthString

	globals.Connection.AddCallback("001", func(e *irc.Event) {
		globals.Connection.SendRaw("CAP REQ :twitch.tv/membership")
		globals.Connection.Join(channel)
	})
	globals.Connection.AddCallback("366", func(e *irc.Event) {})
	globals.Connection.AddCallback("PART", func(e *irc.Event) {
		nick := strings.ToLower(e.Nick)
		if nick == strings.ToLower(botnick) {
			return
		}
		err := bolt.CreateOrUpdateUser(structs.User{
			Name:     nick,
			LastPart: time.Now(),
		})
		if err != nil {
			panic(err)
		}
	})
	globals.Connection.AddCallback("JOIN", func(e *irc.Event) {
		nick := strings.ToLower(e.Nick)
		if nick == strings.ToLower(botnick) {
			return
		}
		err := bolt.CreateOrUpdateUser(structs.User{
			Name:     nick,
			LastJoin: time.Now(),
		})
		if err != nil {
			panic(err)
		}
	})

	globals.Connection.AddCallback("PRIVMSG", func(e *irc.Event) {
		nick := strings.ToLower(e.Nick)
		if nick == strings.ToLower(botnick) {
			return
		}
		err := bolt.CreateOrUpdateUser(structs.User{
			Name:       nick,
			LastActive: time.Now(),
		})
		if err != nil {
			panic(err)
		}
		message := e.Message()
		if len(message) > 1 && strings.HasPrefix(message, "!") {
			i := strings.Index(message, " ")
			channel := e.Arguments[0]
			sender := nick
			var command string
			var params string
			if i < 0 {
				command = message[1:]
				params = ""
			} else {
				command = message[1:i]
				params = message[i:]
			}
			if _, ok := handlers.CommandMap[command]; ok {
				queue.AddJob(globals.CommandQueueName, structs.Job{
					Command: command,
					Channel: channel,
					Sender:  sender,
					Message: message,
					Params: params,
				})
			}
		}
	})

	go queue.HandleCommand(queue.JobChannels[globals.CommandQueueName])

	err = globals.Connection.Connect(serverSSL)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	globals.Connection.Loop()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
