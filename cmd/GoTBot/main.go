package main

import (
	"github.com/thoj/go-ircevent"
	"fmt"
	"strings"
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/joho/godotenv"
	"log"
	"github.com/3stadt/GoTBot/src/queue"
	"github.com/3stadt/GoTBot/src/handlers"
	"github.com/3stadt/GoTBot/src/context"
	"strconv"
	"io/ioutil"
	"os"
	"github.com/BurntSushi/toml"
	"github.com/3stadt/GoTBot/src/db"
	"github.com/3stadt/GoTBot/src/twitch"
)

const serverSSL = "irc.chat.twitch.tv:443"

func main() {
	var err error
	_ = initPlugins()
	db.Up()
	defer db.Down()
	context.Conf, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	queue.NewQueue(context.CommandQueueName)
	botNick := context.Conf["TWITCH_USER"]
	oauth := context.Conf["OAUTH"]
	tw := twitch.Init(oauth, botNick, context.CommandQueueName)
	debug, debugErr := strconv.ParseBool(context.Conf["DEBUG"])
	if debugErr != nil {
		debug = false
	}
	checkErr(err)
	oauthString := strings.TrimSpace(string(oauth))
	tw.Connection.VerboseCallbackHandler = debug
	tw.Connection.Debug = debug
	tw.Connection.UseTLS = true
	tw.Connection.Password = oauthString
	tw.Connection.AddCallback("001", tw.Connected)
	tw.Connection.AddCallback("366", func(e *irc.Event) {})
	tw.Connection.AddCallback("NOTICE", tw.Notice)
	tw.Connection.AddCallback("PART", tw.Part)
	tw.Connection.AddCallback("JOIN", tw.Join)
	tw.Connection.AddCallback("PRIVMSG", tw.Privmsg)
	go queue.HandleCommand(queue.JobQueue[context.CommandQueueName], tw.Connection)
	if err = tw.Connection.Connect(serverSSL); err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	tw.Connection.Loop()
}
func initPlugins() (error) {
	files, err := ioutil.ReadDir("./custom/plugins")
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			if _, err := os.Stat("./custom/plugins/" + file.Name() + "/config.toml"); !os.IsNotExist(err) {
				tomlData, err := ioutil.ReadFile("./custom/plugins/" + file.Name() + "/config.toml")
				if err != nil {
					log.Fatal(err)
					continue
				}
				var commands structs.Commands
				if _, err := toml.Decode(string(tomlData), &commands); err != nil {
					log.Fatal(err)
				}
				for _, c := range commands.Command {
					if _, ok := handlers.PluginCommandMap[c.Name]; ok {
						handlers.PluginCommandMap[c.Name] = append(handlers.PluginCommandMap[c.Name], "./custom/plugins/"+file.Name()+"/"+c.EntryScript)
					} else {
						handlers.PluginCommandMap[c.Name] = []string{"./custom/plugins/" + file.Name() + "/" + c.EntryScript}
					}
				}
			}
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
