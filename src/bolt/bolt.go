package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/3stadt/GoTBot/src/structs"
	"encoding/json"
	"github.com/imdario/mergo"
)

var db *bolt.DB
var userbucket = "Users"

func CreateOrUpdateUser(updateUser structs.User) error {
	baseUser := GetUser(updateUser.Name)
	if baseUser != nil {
		if err := mergo.MergeWithOverwrite(&baseUser, updateUser); err != nil {
			panic(err)
		}
	} else {
		baseUser = &updateUser
	}
	open()
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(userbucket))
		err = b.Put([]byte(baseUser.Name), marshalUser(*baseUser))
		return err
	})
}

func GetUser(username string) *structs.User {
	open()
	defer db.Close()
	var v []byte = nil
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userbucket))
		if b == nil {
			return nil
		}
		v = b.Get([]byte(username))
		return nil
	})
	if v == nil {
		return nil
	}
	user := unmarshalUser(v)
	return &user
}

func marshalUser(user structs.User) []byte {
	jUser, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	return jUser
}

func unmarshalUser(bytes []byte) structs.User {
	user := structs.User{}
	err := json.Unmarshal(bytes, user)
	if err != nil {
		panic(err)
	}
	return user
}

func open(){
	var err error
	db, err = bolt.Open("gotbot.db", 0600, nil)
	if err != nil {
		panic(err)
	}
}