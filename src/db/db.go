package db

import (
	"github.com/asdine/storm"
	"github.com/3stadt/GoTBot/src/context"
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/imdario/mergo"
	"github.com/3stadt/GoTBot/src/errors"
	"time"
	"fmt"
)

func Up() {
	var err error
	context.DB, err = storm.Open(context.DbFile)
	if err != nil {
		panic(err)
	}
	context.Users = context.DB.From("users")
}

func Down() {
	context.DB.Close()
}

func UpdateUser(user structs.User) error {
	baseUser := structs.User{}
	err := context.Users.Get(context.UserBucketName, user.Name, &baseUser)
	if err != nil {
		return &fail.NoTargetUser{Name: user.Name}
	}
	mergo.Merge(&baseUser, user)
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
	mergeUserData(&baseUser, &user)
	return SetUser(baseUser)
}

func mergeUserData(baseUser *structs.User, user *structs.User) {

	// Fix time values not being recognized as empty
	baseUser = clearUser(baseUser)
	user = clearUser(user)

	mergo.MergeWithOverwrite(&baseUser, user)
}

func clearUser(user *structs.User) *structs.User {
	if user.LastActive.IsZero() {
		user.LastActive = nil
		fmt.Println("PING")
	}
	fmt.Println(user.LastActive)
	if user.LastJoin.IsZero() {
		user.LastJoin = nil
		fmt.Println("PONG")
	}
	if user.LastPart.IsZero() {
		user.LastPart = nil
		fmt.Println("PANG")
	}
	if user.FirstSeen.IsZero() {
		user.FirstSeen = nil
		fmt.Println("PÄNG")
	}
	return user
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