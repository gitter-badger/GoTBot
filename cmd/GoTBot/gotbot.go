package GoTBot

import (
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

func Run() {
	var err error
	_ = initPlugins()
	db.Up()
	defer db.Down()
	context.Conf, err = godotenv.Read()
	checkErr(err)
	_, err = connectToTwitch()
	checkErr(err)
}
func connectToTwitch() (twitch.Client, error) {
	botNick := context.Conf["TWITCH_USER"]
	oauth := context.Conf["OAUTH"]
	debug, debugErr := strconv.ParseBool(context.Conf["DEBUG"])
	if debugErr != nil {
		debug = false
	}
	tw := twitch.Init(oauth, botNick, context.CommandQueueName, debug)
	if err := tw.Connection.Connect(serverSSL); err != nil {
		return twitch.Client{}, err
	}
	queue.NewQueue(context.CommandQueueName)
	go queue.HandleCommand(queue.JobQueue[context.CommandQueueName], tw.Connection)
	tw.Connect()
	return tw, nil
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
