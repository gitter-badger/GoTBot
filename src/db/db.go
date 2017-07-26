package db

import (
	"github.com/asdine/storm"
	"github.com/3stadt/GoTBot/src/context"
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/imdario/mergo"
	"github.com/3stadt/GoTBot/src/errors"
	"time"
)

func Up() {
	var err error
	context.DB, err = storm.Open(context.DbFile)
	if err != nil {
		panic(err)
	}
	context.PluginDB, err = storm.Open(context.PluginDbFile)
	if err != nil {
		panic(err)
	}
	context.Users = context.DB.From("users")
}

func Down() {
	context.DB.Close()
	context.PluginDB.Close()
}

func UpdateUser(user structs.User) error {
	baseUser := structs.User{}
	err := context.Users.Get(context.UserBucketName, user.Name, &baseUser)
	if err != nil {
		return &fail.NoTargetUser{Name: user.Name}
	}
	mergo.MergeWithOverwrite(&baseUser, user)
	return SetUser(baseUser)
}

func UpdateMessageCount(nick string) error {
	user := structs.User{}
	err := context.Users.Get(context.UserBucketName, nick, &user)
	if err != nil {
		return CreateUser(nick)
	}
	now := time.Now()
	user.MessageCount++
	user.LastActive = &now
	SetUser(user)
	return nil
}

func CreateUser(nick string) error {
	now := time.Now()
	newUser := structs.User{}
	newUser.Name = nick
	newUser.MessageCount = 1
	newUser.FirstSeen = &now
	newUser.LastActive = &now
	newUser.LastJoin = &now
	return SetUser(newUser)
}

func CreateOrUpdateUser(user structs.User) error {
	baseUser := structs.User{}
	err := context.Users.Get(context.UserBucketName, user.Name, &baseUser)
	if err != nil {
		return SetUser(user)
	}
	mergo.MergeWithOverwrite(&baseUser, user)
	return SetUser(baseUser)
}

func GetUser(name string) (*structs.User, error)  {
	user := structs.User{}
	err := context.Users.Get(context.UserBucketName, name, &user)
	return &user, err
}

func SetUser(user structs.User) error  {
	err := context.Users.Set(context.UserBucketName, user.Name, user)
	return err
}