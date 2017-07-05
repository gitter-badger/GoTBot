package main

import (
	"github.com/thoj/go-ircevent"
	"fmt"
	"strings"
	"github.com/3stadt/GoTBot/structs"
	"github.com/3stadt/GoTBot/bolt"
	"time"
	"github.com/joho/godotenv"
	"log"
	"github.com/3stadt/GoTBot/queue"
	"github.com/3stadt/GoTBot/handlers"
	"github.com/3stadt/GoTBot/globals"
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
	ircConnection := irc.IRC(botnick, botnick)
	ircConnection.VerboseCallbackHandler = debug
	ircConnection.Debug = debug
	ircConnection.UseTLS = true
	ircConnection.Password = oauthString

	ircConnection.AddCallback("001", func(e *irc.Event) {
		ircConnection.SendRaw("CAP REQ :twitch.tv/membership")
		ircConnection.Join(channel)
	})

	ircConnection.AddCallback("366", func(e *irc.Event) {})

	ircConnection.AddCallback("PART", func(e *irc.Event) {
		nick := strings.ToLower(e.Nick)
		if nick == strings.ToLower(botnick) {
			return
		}
		now := time.Now()
		err := bolt.CreateOrUpdateUser(structs.User{
			Name:     nick,
			LastPart: &now,
		})
		if err != nil {
			panic(err)
		}
	})
	ircConnection.AddCallback("JOIN", func(e *irc.Event) {
		nick := strings.ToLower(e.Nick)
		if nick == strings.ToLower(botnick) {
			return
		}
		now := time.Now()
		err := bolt.CreateOrUpdateUser(structs.User{
			Name:     nick,
			LastJoin: &now,
		})
		if err != nil {
			panic(err)
		}
	})

	ircConnection.AddCallback("PRIVMSG", func(e *irc.Event) {
		nick := strings.ToLower(e.Nick)
		now := time.Now()
		if err := bolt.CreateOrUpdateUser(structs.User{
			Name:       nick,
			LastActive: &now,
		}); err != nil {
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

	go queue.HandleCommand(queue.JobQueue[globals.CommandQueueName], ircConnection)

	if err = ircConnection.Connect(serverSSL); err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	ircConnection.Loop()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
