package db

import (
	"github.com/asdine/storm"
	"github.com/3stadt/GoTBot/src/context"
	"github.com/3stadt/GoTBot/src/structs"
	"github.com/imdario/mergo"
	"github.com/3stadt/GoTBot/src/errors"
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

func CreateOrUpdateUser(user structs.User) error {
	baseUser := structs.User{}
	err := context.Users.Get(context.UserBucketName, user.Name, &baseUser)
	if err != nil {
		return SetUser(user)
	}
	mergo.Merge(&baseUser, user)
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