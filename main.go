package main

import (
	"github.com/thoj/go-ircevent"
	"crypto/tls"
	"fmt"
	"strings"
	"math/rand"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
)

const channel = "#3stadt"
const serverSSL = "irc.chat.twitch.tv:443"
const botnick = "3stadt"
const sqliteFileName = "data.sqlite"

var sqlite *sql.DB
var sqliteError error

var commandMap = map[string]func(channel string, sender string, params string, connection *irc.Connection){
	"goSay":   echo,
	"slap":    slap,
	"klatsch": slap,
}

type User struct {
	id         int
	name       string
	lastJoin   time.Time
	lastPart   time.Time
	lastActive time.Time
	firstSeen  time.Time
}

func main() {
	oauth, err := ioutil.ReadFile(".oauth") // just pass the file name
	checkErr(err)
	oauthString := strings.TrimSpace(string(oauth))
	initDB()
	connection := irc.IRC(botnick, botnick)
	connection.VerboseCallbackHandler = false
	connection.Debug = false
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
		openDB()
		logPart(e.Nick)
		closeDB()
	})
	connection.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick == botnick {
			return
		}
		openDB()
		logJoin(e.Nick)
		closeDB()
	})

	connection.AddCallback("PRIVMSG", func(e *irc.Event) {
		if e.Nick == botnick {
			return
		}
		openDB()
		logActive(e.Nick)
		closeDB()
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
				openDB()
				user := getUserData(sender)
				if user.id == 0 {
					initUser(sender)
				}
				commandMap[command](channel, sender, params, connection)
				closeDB()
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

func logActive(username string) {
	updateUserTime(username, "last_part")
}

func logPart(username string) {
	updateUserTime(username, "last_active")
}

func logJoin(username string) {
	updateUserTime(username, "last_join")
}

func updateUserTime(username string, field string) {
	username = strings.TrimSpace(username)
	user := getUserData(username)
	if user.id == 0 {
		initUser(username)
		return
	}
	tx, err := sqlite.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("UPDATE `userinfo` SET `" + field + "` =  datetime('now') WHERE `username`=?")
	checkErr(err)
	_, err = stmt.Exec(username)
	checkErr(err)
	tx.Commit()
	stmt.Close()
}

func initUser(username string) User {
	username = strings.TrimSpace(username)
	tx, err := sqlite.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("INSERT INTO `userinfo` (username, last_join, last_active, first_seen) VALUES (?, datetime('now'), datetime('now'), datetime('now'))")
	checkErr(err)
	_, err = stmt.Exec(username)
	checkErr(err)
	tx.Commit()
	stmt.Close()
	return getUserData(username)
}

func getUserData(username string) User {
	username = strings.TrimSpace(username)
	stmt, err := sqlite.Prepare("SELECT uid, last_join, last_part, last_active, first_seen FROM `userinfo` WHERE `username`=?")
	checkErr(err)
	u := User{
		name: username,
	}
	row := stmt.QueryRow(username)
	row.Scan(&u.id, &u.lastJoin, &u.lastPart, &u.lastActive, &u.firstSeen)
	stmt.Close()
	return u
}

func initDB() {
	sqlStmt, err := ioutil.ReadFile("table_userinfo.sql") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	openDB()
	_, err = sqlite.Exec(string(sqlStmt))
	checkErr(err)
	closeDB()
}

func echo(channel string, sender string, params string, connection *irc.Connection) {
	connection.Privmsg(channel, "Thanks's for sending goSay with '"+params+"' on "+channel+", "+sender+"!")
}

func slap(channel string, sender string, params string, connection *irc.Connection) {
	victim := strings.TrimSpace(params)
	if len(params) < 1 || strings.ContainsAny(victim, " ") {
		return
	}

	if victim == "himself" || victim == "herself" || victim == "itself" {
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

func openDB() {
	sqlite, sqliteError = sql.Open("sqlite3", "./" + sqliteFileName)
	checkErr(sqliteError)
}

func closeDB() {
	if sqlite != nil {
		sqlite.Close()
	}
}

func checkErr(err error) {
	if err != nil {
		closeDB()
		panic(err)
	}
}
