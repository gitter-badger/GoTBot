package GoTBot

import (
	"github.com/3stadt/GoTBot/src/structs"
	"log"
	"github.com/3stadt/GoTBot/src/queue"
	"github.com/3stadt/GoTBot/src/handlers"
	"strconv"
	"io/ioutil"
	"os"
	"github.com/BurntSushi/toml"
	"github.com/3stadt/GoTBot/src/db"
	"github.com/3stadt/GoTBot/src/twitch"
	"github.com/3stadt/GoTBot/src/res"
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"github.com/joho/godotenv"
)

const serverSSL = "irc.chat.twitch.tv:443"

func Run() {
	_ = initPlugins()
	p := &db.Pool{
		DbFile:       "gotbot.db",
		PluginDbFile: "gotbotPlugins.db",
	}
	p.Up()
	defer p.Down()

	r := gin.Default()
	
	r.Delims("{[{", "}]}")
	r.LoadHTMLGlob("assets/templates/*")

	r.GET("/sampleMessage", func(c *gin.Context) {
		t := time.Now()
		c.JSON(200, gin.H{
			"message": "Hello Twitch! The button was last clicked at " + t.Format("15:04:05"),
		})
	})

	r.GET("/twitch/disconnect", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Not yet implemented",
		})
	})

	r.POST("/twitch/connect", func(c *gin.Context) {
		cfg, err := godotenv.Read()
		rs := &res.Vars{
			Conf:      cfg,
			Constants: res.GetConst(),
		}
		checkErr(err)
		tw, err := connectToTwitch(p, rs)
		checkErr(err)
		go tw.Connection.Loop()
		c.JSON(200, gin.H{
			"message": "Connected",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "GoTBot",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func connectToTwitch(p *db.Pool, rs *res.Vars) (twitch.Client, error) {
	botNick := rs.Conf["TWITCH_USER"]
	oauth := rs.Conf["OAUTH"]
	debug, debugErr := strconv.ParseBool(rs.Conf["DEBUG"])
	if debugErr != nil {
		debug = false
	}
	tw := twitch.Init(oauth, botNick, rs.Constants.CommandQueueName, debug, p, rs)
	if err := tw.Connection.Connect(serverSSL); err != nil {
		return twitch.Client{}, err
	}
	queue.NewQueue(rs.Constants.CommandQueueName)
	go queue.HandleCommand(queue.JobQueue[rs.Constants.CommandQueueName], tw.Connection, p, rs)
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
